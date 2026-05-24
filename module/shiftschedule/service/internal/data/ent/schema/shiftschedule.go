// ent/schema/employeeshiftassignment.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// EmployeeShiftSchedule
// EmployeeShiftAssignment holds the schema definition for the EmployeeShiftAssignment entity.
type ShiftSchedule struct {
	ent.Schema
}

// Fields of the EmployeeShiftAssignment.
func (ShiftSchedule) Fields() []ent.Field {
	return []ent.Field{
		// field.UUID("id", uuid.UUID{}).
		// 	Default(uuid.New).
		// 	Unique(),
		// field.UUID("id", uuid.UUID{}).Default(uuid.New),
		// field.UUID("employee_id", uuid.UUID{}).Immutable(), // Foreign key logical to employers service
		field.Int64("id"),
		field.String("schedule_code"),
		field.String("karya_code"),
		// field.Time("tanggal").Immutable(), // DATE in PostgreSQL
		field.Time("tanggal").
			SchemaType(map[string]string{
				"postgres": "DATE",
			}),
		field.String("created_by"),
		field.String("depart_code"),
		field.String("shift_id"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// // Edges of the EmployeeShiftAssignment.
// func (EmployeeShiftAssignment) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.From("shift", Shift.Type).Ref("assignments").Unique().Required(),
// 	}
// }
