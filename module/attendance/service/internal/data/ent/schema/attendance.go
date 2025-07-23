// ent/schema/attendance.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Attendance holds the schema definition for the Attendance entity.
type Attendance struct {
	ent.Schema
}

// Fields of the Attendance.
func (Attendance) Fields() []ent.Field {
	return []ent.Field{
		// field.UUID("id", uuid.UUID{}).
		// 	Default(uuid.New).
		// 	Unique(),
		// field.Int64("id"),
		// field.UUID("id", uuid.UUID{}).Default(uuid.New),
		// field.UUID("employee_id", uuid.UUID{}).Immutable(),                 // Foreign key logical to employers service
		field.Int64("id"),
		field.String("karyacode"),
		// field.Time("clock_in_time").Immutable(),            // TIMESTAMPTZ in PostgreSQL
		field.Time("clock_in_time"),                        // TIMESTAMPTZ in PostgreSQL
		field.Time("clock_out_time").Optional().Nillable(), // TIMESTAMPTZ in PostgreSQL, can be NULL
		// field.UUID("assigned_shift_id", uuid.UUID{}).Optional().Nillable(), // Foreign key to Shift
		field.String("status").MaxLen(50).Optional().Nillable(),
		field.Int("effective_work_duration_minutes").Optional().Nillable(),
		field.Int("overtime_minutes").Default(0),
		field.Text("notes").Optional().Nillable(),
		field.String("location").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// // Edges of the Attendance.
// func (Attendance) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.To("overtime", Overtime.Type),
// 		edge.From("assigned_shift", Shift.Type).Ref("attendances").Unique().Field("assigned_shift_id"),
// 	}
// }
