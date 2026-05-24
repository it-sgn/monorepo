package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "mall-go/api/organization/service/v1"
	"mall-go/module/organization/service/internal/biz"
	"mall-go/module/organization/service/internal/data/model/assignment"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ biz.AssignmentRepo = (*assignmentRepo)(nil)

type assignmentRepo struct {
	data        *Data
	Log         *log.Helper
	redisClient *redis.Client
}

func NewAssignmentRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.AssignmentRepo {
	return &assignmentRepo{
		data:        data,
		Log:         log.NewHelper(log.With(logger, "module", "data/assignment")),
		redisClient: rdb,
	}
}

func (r *assignmentRepo) AssignPosition(ctx context.Context, req *biz.AssignmentData) (*v1.AssignmentResponse, error) {
	if req.EmployeeID == "" || req.PositionID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: employee_id and position_id are required")
	}
	// layout := "02-01-2006 15:04"
	// parsedTime, err := time.Parse(layout, req.StartDate)
	// if err != nil {
	// 	log.Fatalf("Gagal parse waktu: %v", err)
	// }
	// Gunakan langsung req.StartDate (sudah time.Time)
	po, err := r.data.db.Assignment.
		Create().
		SetEmployeeID(req.EmployeeID).
		SetPositionID(req.PositionID).
		SetStartDate(req.StartDate). // ✅ simpan langsung time.Time ke Ent
		Save(ctx)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to assign position: %v", err)
	}

	// Convert end_date jika ada
	// var endDate *timestamppb.Timestamp
	// if po.EndDate != nil {
	// 	endDate = timestamppb.New(*po.EndDate)
	// }

	// Build response
	resp := &v1.AssignmentResponse{
		Assignment: &v1.Assignment{
			Id:         po.ID,
			EmployeeId: po.EmployeeID,
			PositionId: po.PositionID,
			StartDate:  po.StartDate.Format(time.RFC3339),
			EndDate:    po.EndDate.Format(time.RFC3339),
			CreatedAt:  time.Now().Format(time.RFC3339),
		},
	}

	// Cache ke Redis
	cacheKey := "assignment:" + po.EmployeeID
	if bytes, err := json.Marshal(resp); err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, bytes, 5*time.Minute).Err()
	}

	return resp, nil
}

// func (r *assignmentRepo) AssignPosition(ctx context.Context, req *biz.AssignmentData) (*v1.AssignmentResponse, error) {
// 	if req.EmployeeID == "" || req.PositionID <= 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "invalid request: employee_id and position_id are required")
// 	}

// 	po, err := r.data.db.Assignment.
// 		Create().
// 		SetEmployeeID(req.EmployeeID).
// 		SetPositionID(req.PositionID).
// 		SetStartDate(req.StartDate).
// 		Save(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to assign position: %v", err)
// 	}

// 	var endDate *timestamppb.Timestamp
// 	if po.EndDate != nil {
// 		endDate = timestamppb.New(*po.EndDate)
// 	}

// 	resp := &v1.AssignmentResponse{
// 		Assignment: &v1.Assignment{
// 			Id:         po.ID,
// 			EmployeeId: po.EmployeeID,
// 			PositionId: strconv.FormatInt(po.PositionID, 10),
// 			StartDate:  timestamppb.New(po.StartDate),
// 			EndDate:    endDate,
// 			CreatedAt:  timestamppb.New(po.CreateTime),
// 		},
// 	}

// 	cacheKey := "assignment:" + po.EmployeeID
// 	if bytes, err := json.Marshal(resp); err == nil {
// 		_ = r.redisClient.Set(ctx, cacheKey, bytes, 5*time.Minute).Err()
// 	}

// 	return resp, nil
// }

func (r *assignmentRepo) GetAssignment(ctx context.Context, c *v1.GetAssignmentRequest) (*v1.AssignmentResponse, error) {
	cacheKey := "assignment:" + c.EmployeeId
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached v1.AssignmentResponse
		if unmarshalErr := json.Unmarshal([]byte(val), &cached); unmarshalErr == nil {
			r.Log.Infof("Cache hit: assignment %s", c.EmployeeId)
			return &cached, nil
		}
	}

	po, err := r.data.db.Assignment.Query().
		Where(assignment.EmployeeID(c.EmployeeId)).
		First(ctx)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "assignment not found: %v", err)
	}

	resp := &v1.AssignmentResponse{
		Assignment: &v1.Assignment{
			Id:         int64(po.ID),
			EmployeeId: po.EmployeeID,
			PositionId: po.PositionID,
			StartDate:  po.StartDate.Format(time.RFC3339),
			EndDate:    po.EndDate.Format(time.RFC3339),
			CreatedAt:  po.CreateTime.Format(time.RFC3339),
		},
	}
	// if po.EndDate != nil {
	// 	resp.Assignment.EndDate = po.EndDate.Format(time.RFC3339)
	// }

	if bytes, err := json.Marshal(resp); err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, bytes, 5*time.Minute).Err()
	}

	return resp, nil
}

func (r *assignmentRepo) DeleteAssignment(ctx context.Context, id int64) error {
	po, err := r.data.db.Assignment.Get(ctx, id)
	if err != nil {
		return status.Errorf(codes.NotFound, "assignment not found")
	}

	err = r.data.db.Assignment.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete assignment: %v", err)
	}

	// Invalidate related cache
	_ = r.redisClient.Del(ctx,
		"assignment:"+po.EmployeeID,
		"assignments:list",
		fmt.Sprintf("assignmentDetail:%d", id),
	).Err()

	r.Log.Infof("Assignment %d deleted", id)
	return nil
}

func (r *assignmentRepo) Count(ctx context.Context) (int, error) {
	return r.data.db.Assignment.Query().Count(ctx)
}

func (r *assignmentRepo) ListAssignment(ctx context.Context, pageNum, pageSize int64) (*v1.ListAssignmentsResponse, int, error) {
	cacheKey := fmt.Sprintf("assignments:list:page:%d:size:%d", pageNum, pageSize)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached struct {
			Assignments []*v1.Assignment
			Total       int
		}
		if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
			r.Log.Infof("Cache hit: ListAssignments page=%d", pageNum)
			return &v1.ListAssignmentsResponse{Assignments: cached.Assignments}, cached.Total, nil
		}
	}

	query := r.data.db.Assignment.Query()
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, status.Errorf(codes.Internal, "failed to count assignments: %v", err)
	}

	assignments, err := query.
		Offset(int((pageNum - 1) * pageSize)).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, status.Errorf(codes.Internal, "failed to list assignments: %v", err)
	}

	list := make([]*v1.Assignment, 0, len(assignments))
	for _, a := range assignments {
		item := &v1.Assignment{
			Id:         int64(a.ID),
			EmployeeId: a.EmployeeID,
			PositionId: a.PositionID,
			StartDate:  a.StartDate.Format(time.RFC3339),
			EndDate:    a.EndDate.Format(time.RFC3339),
			CreatedAt:  a.CreateTime.Format(time.RFC3339),
		}
		// if a.EndDate != nil {
		// 	item.EndDate = timestamppb.New(*a.EndDate)
		// }
		list = append(list, item)
	}

	toCache := struct {
		Assignments []*v1.Assignment
		Total       int
	}{
		Assignments: list,
		Total:       total,
	}
	if bytes, err := json.Marshal(toCache); err == nil {
		_ = r.redisClient.Set(ctx, cacheKey, bytes, 3*time.Minute).Err()
	}

	return &v1.ListAssignmentsResponse{Assignments: list}, total, nil
}

// package data

// import (
// 	"context"
// 	"encoding/json" // Digunakan untuk serialisasi/deserialisasi JSON
// 	"fmt"
// 	"strconv"
// 	"time" // Digunakan untuk mengatur waktu kedaluwarsa cache

// 	"mall-go/module/organization/service/internal/biz"
// 	"mall-go/module/organization/service/internal/data/model/assignment"
// 	"mall-go/pkg/utils/pagination"

// 	v1 "mall-go/api/organization/service/v1"

// 	"github.com/go-kratos/kratos/v2/log"
// 	"github.com/go-redis/redis/v8" // Pastikan versi Redis client yang benar
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// // Pastikan antarmuka AssignmentRepo diimpor dengan benar dari biz package
// var _ biz.AssignmentRepo = (*assignmentRepo)(nil)

// type assignmentRepo struct {
// 	data        *Data
// 	Log         *log.Helper
// 	redisClient *redis.Client
// }

// func NewAssignmentRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.AssignmentRepo {
// 	return &assignmentRepo{
// 		data:        data,
// 		redisClient: rdb,
// 		Log:         log.NewHelper(log.With(logger, "module", "data/assignment")),
// 	}
// }
// func (r *assignmentRepo) GetAssignment(ctx context.Context, c *v1.GetAssignmentRequest) (*v1.AssignmentResponse, error) {
// 	cacheKey := "assignment:" + c.EmployeeId
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedAssignment v1.AssignmentResponse
// 		if unmarshalErr := json.Unmarshal([]byte(val), &cachedAssignment); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for assignment ID: %s", c.EmployeeId)
// 			return &cachedAssignment, nil
// 		}
// 		r.Log.Warnf("Failed to unmarshal cached assignment for ID %s: %v", c.EmployeeId, err)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for assignment ID %s: %v", c.EmployeeId, err)
// 	}

// 	// Query dari DB
// 	po, err := r.data.db.Assignment.Query().
// 		Where(assignment.EmployeeID(c.EmployeeId)).First(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	assignmentBiz := &v1.AssignmentResponse{
// 		Assignment: &v1.Assignment{
// 			Id:         int64(po.ID),
// 			EmployeeId: po.EmployeeID,
// 			PositionId: string(po.PositionID),
// 			StartDate:  timestamppb.New(po.StartDate),
// 			EndDate:    timestamppb.New(*po.EndDate), // Jika nullable, cek null
// 			CreatedAt:  timestamppb.New(po.CreateTime),
// 		},
// 	}

// 	// Simpan cache
// 	if assignmentBytes, marshalErr := json.Marshal(assignmentBiz); marshalErr == nil {
// 		r.redisClient.Set(ctx, cacheKey, assignmentBytes, 3*time.Minute)
// 	} else {
// 		r.Log.Errorf("Failed to marshal assignment for caching: %v", marshalErr)
// 	}

// 	return assignmentBiz, nil
// }
// func (r *assignmentRepo) AssignPosition(ctx context.Context, req *biz.AssignmentData) (*v1.AssignmentResponse, error) {
// 	// Validasi input
// 	if req.EmployeeID == "" {
// 		return nil, status.Errorf(codes.InvalidArgument, "employee_id is required")
// 	}
// 	if req.PositionID <= 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "position_id must be greater than 0")
// 	}
// 	if req.StartDate.IsZero() {
// 		return nil, status.Errorf(codes.InvalidArgument, "start_date is required")
// 	}

// 	// Simpan ke database
// 	po, err := r.data.db.Assignment.
// 		Create().
// 		SetEmployeeID(req.EmployeeID).
// 		SetPositionID(req.PositionID).
// 		SetStartDate(req.StartDate).
// 		Save(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create assignment: %v", err)
// 	}

// 	// Siapkan EndDate jika ada
// 	var endDate *timestamppb.Timestamp
// 	if po.EndDate != nil {
// 		endDate = timestamppb.New(*po.EndDate)
// 	}

// 	// Build response sesuai dengan .proto (v1.Assignment)
// 	resp := &v1.AssignmentResponse{
// 		Assignment: &v1.Assignment{
// 			Id:         po.ID,
// 			EmployeeId: po.EmployeeID,
// 			PositionId: strconv.FormatInt(po.PositionID, 10),
// 			StartDate:  timestamppb.New(po.StartDate),
// 			EndDate:    endDate,
// 			CreatedAt:  timestamppb.New(po.CreateTime),
// 		},
// 	}

// 	// Simpan ke Redis cache (opsional)
// 	cacheKey := "assignment:" + po.EmployeeID
// 	if assignmentBytes, err := json.Marshal(resp); err == nil {
// 		if err := r.redisClient.Set(ctx, cacheKey, assignmentBytes, 5*time.Minute).Err(); err != nil {
// 			r.Log.Warnf("failed to cache assignment: %v", err)
// 		}
// 	} else {
// 		r.Log.Warnf("failed to marshal assignment response: %v", err)
// 	}

// 	return resp, nil
// }

// // func (r *assignmentRepo) AssignPosition(ctx context.Context, req *biz.AssignmentData) (*biz.AssignmentData, error) {
// // 	// Validasi input (disarankan)
// // 	if req.EmployeeID == "" || req.PositionID <= 0 || req.StartDate.IsZero() {
// // 		return nil, status.Errorf(codes.InvalidArgument, "employee_id, position_id, and start_date are required")
// // 	}
// // 	// positionid, _ := strconv.ParseInt(req.PositionID, 10, 64)

// // 	// Simpan ke database
// // 	po, err := r.data.db.Assignment.
// // 		Create().
// // 		SetEmployeeID(req.EmployeeID).
// // 		SetPositionID(req.PositionID).
// // 		SetStartDate(req.StartDate).
// // 		Save(ctx)
// // 	if err != nil {
// // 		return nil, status.Errorf(codes.Internal, "failed to create assignment: %v", err)
// // 	}

// // 	// Build response
// // 	resp := &biz.AssignmentData{
// // 		ID:         po.ID,
// // 		EmployeeID: po.EmployeeID,
// // 		PositionID: po.PositionID,
// // 		StartDate:  po.StartDate,
// // 		EndDate:    *po.EndDate,
// // 		CreatedAt:  po.CreateTime,
// // 		UpdatedAt:  po.UpdateTime,
// // 		CreatedBy:  po.CreatedBy,
// // 		UpdatedBy:  po.UpdatedBy,
// // 	}

// // 	// Optional: Simpan ke Redis cache
// // 	cacheKey := "assignment:" + po.EmployeeID
// // 	if assignmentBytes, err := json.Marshal(resp); err == nil {
// // 		if err := r.redisClient.Set(ctx, cacheKey, assignmentBytes, 5*time.Minute).Err(); err != nil {
// // 			r.Log.Warnf("failed to cache assignment: %v", err)
// // 		}
// // 	}

// // 	return resp, nil
// // }

// func (r assignmentRepo) DeleteAssignment(ctx context.Context, id int64) error {
// 	// Dapatkan detail departemen sebelum menghapus untuk menginvalidasi cache berdasarkan kode juga
// 	deptToDelete, err := r.data.db.Assignment.Get(ctx, id)
// 	if err != nil {
// 		return err // Atau log error dan tetap coba hapus jika tidak kritis
// 	}

// 	err = r.data.db.Assignment.
// 		DeleteOneID(id).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	// Invalidate cache setelah menghapus
// 	r.redisClient.Del(ctx, "assignments:list")
// 	r.redisClient.Del(ctx, "assignment:"+string(id))
// 	// Asumsi ada cache 'assignmentDetail' jika digunakan di tempat lain
// 	r.redisClient.Del(ctx, "assignmentDetail:"+string(id))
// 	if deptToDelete != nil {
// 		r.redisClient.Del(ctx, "assignment_code:"+deptToDelete.EmployeeID) // Invalidate cache by code
// 	}
// 	r.Log.Infof("Assignment dengan ID %d terhapus", id)
// 	return nil
// }

// func (r *assignmentRepo) Count(ctx context.Context) (int, error) {
// 	// Pertimbangkan caching untuk count jika sering diakses dan tidak perlu sangat up-to-date
// 	dt, _ := r.data.db.Assignment.Query().Count(ctx)
// 	return dt, nil
// }

// func (r *assignmentRepo) ListAssignment(ctx context.Context, pageNum, pageSize int64) (*v1.ListAssignmentsResponse, int, error) {
// 	cacheKey := fmt.Sprintf("assignments:list:page:%d:size:%d", pageNum, pageSize)

// 	// Try get from Redis
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedResponse struct {
// 			Assignments []*v1.Assignment
// 			Total       int
// 		}
// 		var unmarshalErr error
// 		if unmarshalErr = json.Unmarshal([]byte(val), &cachedResponse); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for ListAssignments page %d, size %d", pageNum, pageSize)
// 			return &v1.ListAssignmentsResponse{
// 				Assignments: cachedResponse.Assignments,
// 			}, cachedResponse.Total, nil
// 		}
// 		r.Log.Warnf("Failed to unmarshal cached ListAssignments: %v", unmarshalErr)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for ListAssignments: %v", err)
// 	}

// 	// If not in cache or failed, query DB
// 	query := r.data.db.Assignment.Query()

// 	// Get total count for pagination
// 	total, err := query.Clone().Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Get paginated data
// 	pos, err := query.
// 		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
// 		Limit(int(pageSize)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	assignments := make([]*v1.Assignment, 0, len(pos))
// 	for _, po := range pos {
// 		assignment := &v1.Assignment{
// 			Id:         int64(po.ID),
// 			EmployeeId: po.EmployeeID,
// 			PositionId: strconv.FormatInt(po.PositionID, 10),
// 			StartDate:  timestamppb.New(po.StartDate),
// 			CreatedAt:  timestamppb.New(po.CreateTime),
// 		}
// 		// Handle optional EndDate
// 		if po.EndDate != nil {
// 			assignment.EndDate = timestamppb.New(*po.EndDate)
// 		}
// 		assignments = append(assignments, assignment)
// 	}

// 	// Cache result to Redis
// 	responseToCache := struct {
// 		Assignments []*v1.Assignment
// 		Total       int
// 	}{
// 		Assignments: assignments,
// 		Total:       total,
// 	}

// 	if responseBytes, marshalErr := json.Marshal(responseToCache); marshalErr == nil {
// 		if err := r.redisClient.Set(ctx, cacheKey, responseBytes, 3*time.Minute).Err(); err != nil {
// 			r.Log.Warnf("Failed to cache ListAssignments: %v", err)
// 		}
// 	} else {
// 		r.Log.Errorf("Failed to marshal ListAssignments for caching: %v", marshalErr)
// 	}

// 	return &v1.ListAssignmentsResponse{
// 		Assignments: assignments,
// 	}, total, nil
// }
