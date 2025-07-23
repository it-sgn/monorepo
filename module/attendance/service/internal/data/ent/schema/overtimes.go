// ent/schema/overtime.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Overtime holds the schema definition for the Overtime entity.
type Overtime struct {
	ent.Schema
}

// Fields of the Overtime.
func (Overtime) Fields() []ent.Field {
	return []ent.Field{
		// field.UUID("id", uuid.UUID{}).
		// 	Default(uuid.New).
		// 	Unique(),
		// field.UUID("employee_id", uuid.UUID{}).Immutable(), // Foreign key logical to employers service
		field.String("karyacode"),
		field.Time("overtime_start_time").Immutable(),
		field.Time("overtime_end_time").Immutable(),
		field.Int("duration_minutes"),
		// field.UUID("approved_by", uuid.UUID{}).Optional().Nillable(), // Foreign key logical to employers service
		field.String("approved_by"),
		field.String("status").MaxLen(50).Default("pending"),
		field.Text("notes").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// // Edges of the Overtime.
// func (Overtime) Edges() []ent.Edge {
// 	return []ent.Edge{
// 		edge.From("attendance", Attendance.Type).Ref("overtime").Unique().Required(),
// 	}
// }
