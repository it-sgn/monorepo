// ent/schema/shift.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Shift holds the schema definition for the Shift entity.
type Shift struct {
	ent.Schema
}

// Fields of the Shift.
func (Shift) Fields() []ent.Field {
	return []ent.Field{
		// field.UUID("id", uuid.UUID{}).
		// 	Default(uuid.New).
		// 	Unique(),
		// field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Int64("id"),
		field.String("name").MaxLen(100).Unique(),
		field.Time("start_time"), // TIME in PostgreSQL
		field.Time("end_time"),   // TIME in PostgreSQL
		field.Int("break_duration_minutes").Default(0),
		field.String("created_by"),
		field.String("updated_by"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// // Edges of the Shift.
// func (Shift) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("assignments", EmployeeShiftAssignment.Type),
// 		edge.To("attendances", Attendance.Type), // Jika ingin direct link dari shift ke absensi
// 	}
// }
