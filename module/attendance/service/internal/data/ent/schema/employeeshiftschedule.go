// ent/schema/employeeshiftassignment.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// EmployeeShiftSchedule
// EmployeeShiftAssignment holds the schema definition for the EmployeeShiftAssignment entity.
type EmployeeShiftSchedule struct {
	ent.Schema
}

// Fields of the EmployeeShiftAssignment.
func (EmployeeShiftSchedule) Fields() []ent.Field {
	return []ent.Field{
		// field.UUID("id", uuid.UUID{}).
		// 	Default(uuid.New).
		// 	Unique(),
		// field.UUID("id", uuid.UUID{}).Default(uuid.New),
		// field.UUID("employee_id", uuid.UUID{}).Immutable(), // Foreign key logical to employers service
		field.Int64("id"),
		field.String("karyacode"),
		field.Time("assignment_date").Immutable(), // DATE in PostgreSQL
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
