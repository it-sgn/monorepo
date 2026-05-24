package biz

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/robfig/cron"
	"github.com/robfig/cron/v3" // Correct import for cron library
	// "github.com/go-kratos/beer-shop/pkg/page_token"
)

type CronZK struct {
	ID        int64
	Name      string
	Spec      string
	Command   string
	Enabled   bool
	LastRunAt *time.Time
	NextRunAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// cronLogger is an adapter to use Kratos log.Helper with robfig/cron.
type cronLogger struct {
	log *log.Helper
}

// Info logs routine messages about cron's operation.
func (cl *cronLogger) Info(msg string, keysAndValues ...any) {
	cl.log.Infow(append([]any{"msg", msg}, keysAndValues...)...)
}

// Error logs an error condition.
func (cl *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	cl.log.Errorw(append([]any{"msg", msg, "error", err}, keysAndValues...)...)
}

type CronZKRepo interface {
	CreateCronZK(context.Context, *CronZK) (*CronZK, error)
	GetCronZK(context.Context, int64) (*CronZK, error)
	UpdateCronZK(context.Context, *CronZK) (*CronZK, error)
	DeleteCronZK(context.Context, int64) error
	ListCronZK(context.Context, int64, int64) ([]*CronZK, int64, error)
	ListEnabledCronZKs(context.Context) ([]*CronZK, error)
}

// CronZKUseCase is a CronZK usecase.
type CronZKUseCase struct {
	repo        CronZKRepo
	log         *log.Helper
	cron        *cron.Cron
	runningJobs map[string]struct{} // Keep track of currently running jobs
}

func NewCronZKUseCase(repo CronZKRepo, logger log.Logger) *CronZKUseCase {
	cl := &cronLogger{log: log.NewHelper(log.With(logger, "module", "biz/cron-scheduler"))}
	uc := &CronZKUseCase{
		repo:        repo,
		log:         log.NewHelper(log.With(logger, "module", "biz/cronjob")),
		cron:        cron.New(cron.WithChain(cron.Recover(cl)), cron.WithLogger(cl)), // Restored for modern cron/v3
		runningJobs: make(map[string]struct{}),
	}
	uc.cron.Start() // Start the cron scheduler
	return uc
}

func (uc *CronZKUseCase) Create(ctx context.Context, cj *CronZK) (*CronZK, error) {
	uc.log.WithContext(ctx).Infof("CreateCronZK: %v", cj.Name)
	return uc.repo.CreateCronZK(ctx, cj)
}

func (uc *CronZKUseCase) Update(ctx context.Context, cj *CronZK) (*CronZK, error) {
	uc.log.WithContext(ctx).Infof("UpdateCronZK: %v", cj.Name)
	return uc.repo.UpdateCronZK(ctx, cj)
}

func (uc *CronZKUseCase) Delete(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteCronZK: %s", id)
	return uc.repo.DeleteCronZK(ctx, id)
}

func (uc *CronZKUseCase) Get(ctx context.Context, id int64) (*CronZK, error) {
	uc.log.WithContext(ctx).Infof("GetCronZK: %s", id)
	return uc.repo.GetCronZK(ctx, id)
}

func (uc *CronZKUseCase) List(ctx context.Context, page, pageSize int64) ([]*CronZK, int64, error) {
	uc.log.WithContext(ctx).Infof("ListCronZK: page=%d, pageSize=%d", page, pageSize)
	return uc.repo.ListCronZK(ctx, page, pageSize)
}

// ScheduleCronZKs fetches enabled cron jobs from DB and schedules them.
func (uc *CronZKUseCase) ScheduleCronZKs(ctx context.Context) {
	uc.log.WithContext(ctx).Info("Scheduling cron jobs...")
	jobs, err := uc.repo.ListEnabledCronZKs(ctx)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to list enabled cron jobs: %v", err)
		return
	}

	// Get current entries in the scheduler
	currentEntries := make(map[cron.EntryID]struct{})
	for _, entry := range uc.cron.Entries() {
		currentEntries[entry.ID] = struct{}{}
	}

	// Remove entries from cron that are no longer enabled or exist in DB
	for _, entry := range uc.cron.Entries() {
		found := false
		for _, job := range jobs {
			// In modern cron/v3, Entry.ID is directly accessible.
			// We assume job.ID is a string representation of the EntryID for comparison.
			if fmt.Sprintf("%d", entry.ID) == strconv.FormatInt(job.ID, 10) {
				found = true
				break
			}
		}
		if !found {
			uc.cron.Remove(entry.ID) // Correct usage of Remove with EntryID
			uc.log.WithContext(ctx).Infof("Removed cron job from scheduler: %s", fmt.Sprintf("%d", entry.ID))
		}
	}

	// Schedule or re-schedule jobs from the database
	for _, job := range jobs {
		var existingEntryID cron.EntryID = 0
		for _, entry := range uc.cron.Entries() {
			if fmt.Sprintf("%d", entry.ID) == strconv.FormatInt(job.ID, 10) {
				existingEntryID = entry.ID
				break
			}
		}

		if existingEntryID != 0 {
			uc.cron.Remove(existingEntryID) // Correct usage of Remove with EntryID
			uc.log.WithContext(ctx).Infof("Re-scheduling cron job: %s (old EntryID: %d)", job.Name, existingEntryID)
		}

		// Calculate next run time
		schedule, err := cron.ParseStandard(job.Spec)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("Invalid cron spec for job %s: %v", job.Name, err)
			continue
		}
		nextRun := schedule.Next(time.Now())

		// Update next_run_at in DB
		job.NextRunAt = &nextRun
		_, err = uc.repo.UpdateCronZK(ctx, job)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("Failed to update next_run_at for job %s: %v", job.Name, err)
		}

		// Add the job to the cron scheduler
		id, err := uc.cron.AddJob(job.Spec, cron.FuncJob(func(jobID string) func() {
			return func() {
				// uc.executeCronZK(jobID)
				jobID64, _ := strconv.ParseInt(jobID, 10, 64)
				uc.executeCronZK(jobID64)
			}
		}(strconv.FormatInt(job.ID, 10)))) // Pass job.ID to the closure
		if err != nil {
			uc.log.WithContext(ctx).Errorf("Failed to add cron job %s to scheduler: %v", job.Name, err)
			continue
		}
		uc.log.WithContext(ctx).Infof("Scheduled cron job %s (ID: %s, Assigned Cron EntryID: %d)", job.Name, job.ID, id)
	}
}

// executeCronZK executes the command associated with a cron job.
func (uc *CronZKUseCase) executeCronZK(jobID int64) {
	ctx := context.Background()
	uc.log.WithContext(ctx).Infof("Executing cron job: %s", jobID)

	jobID64 := strconv.FormatInt(jobID, 10)
	// Prevent concurrent execution of the same job
	_, running := uc.runningJobs[jobID64]
	if running {
		uc.log.WithContext(ctx).Warnf("Cron job %s is already running, skipping this execution.", jobID)
		return
	}

	uc.runningJobs[jobID64] = struct{}{}
	defer delete(uc.runningJobs, jobID64)

	job, err := uc.repo.GetCronZK(ctx, jobID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to get cron job %s for execution: %v", jobID, err)
		return
	}

	if !job.Enabled {
		uc.log.WithContext(ctx).Infof("Cron job %s is disabled, skipping execution.", jobID)
		return
	}

	cmdParts := splitCommand(job.Command)
	if len(cmdParts) == 0 {
		uc.log.WithContext(ctx).Errorf("Cron job %s has no command to execute.", jobID)
		return
	}

	cmdCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, cmdParts[0], cmdParts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to execute command for job %s: %v\nOutput: %s", job.Name, err, string(output))
	} else {
		uc.log.WithContext(ctx).Infof("Successfully executed command for job %s.\nOutput: %s", job.Name, string(output))
	}

	now := time.Now()
	job.LastRunAt = &now

	schedule, err := cron.ParseStandard(job.Spec)
	if err == nil {
		nextRun := schedule.Next(now)
		job.NextRunAt = &nextRun
	} else {
		uc.log.WithContext(ctx).Errorf("Failed to parse cron spec for job %s after execution: %v", job.Name, err)
	}

	_, err = uc.repo.UpdateCronZK(ctx, job)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to update last_run_at/next_run_at for job %s: %v", job.Name, err)
	}
}

// splitCommand splits a command string into parts, handling quoted arguments.
func splitCommand(command string) []string {
	var parts []string
	var currentPart string
	inQuote := false
	for i := 0; i < len(command); i++ {
		char := command[i]
		if char == '"' {
			inQuote = !inQuote
			if !inQuote && currentPart != "" {
				parts = append(parts, currentPart)
				currentPart = ""
			}
			continue
		}
		if char == ' ' && !inQuote {
			if currentPart != "" {
				parts = append(parts, currentPart)
				currentPart = ""
			}
			continue
		}
		currentPart += string(char)
	}
	if currentPart != "" {
		parts = append(parts, currentPart)
	}
	return parts
}

// StartScheduler starts a goroutine that periodically reloads and schedules cron jobs.
func (uc *CronZKUseCase) StartScheduler(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Initial scheduling
		uc.ScheduleCronZKs(ctx)

		for range ticker.C {
			select {
			case <-ctx.Done():
				uc.log.Info("Scheduler context cancelled, stopping.")
				return
			default:
				uc.ScheduleCronZKs(ctx)
			}
		}
	}()
}

// StopScheduler stops the internal cron scheduler.
func (uc *CronZKUseCase) StopScheduler() {
	uc.cron.Stop()
	uc.log.Info("Cron scheduler stopped.")
}
