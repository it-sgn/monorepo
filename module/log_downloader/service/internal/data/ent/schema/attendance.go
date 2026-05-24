// ent/schema/attendance.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Attendance holds the schema definition for the Attendance entity.
type Attendance struct {
	ent.Schema
}

func (Attendance) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "attendance_log"},
	}
}
func (Attendance) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("Auto increment ID"),
		field.Int("user_id").
			Comment("ID pengguna").
			Positive(),
		field.String("device_ip").
			Comment("IP address fingerprint device").
			NotEmpty(),
		field.Time("att_log").
			Comment("Waktu absensi (timestamp)").
			// Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}),
		field.Int("status").
			NonNegative().
			Max(1),
		// field.Int("status").
		// 	Comment("0 = clock_in, 1 = clock_out"),
		// Positive(),
		// Virtual field untuk date part (harus manual ALTER TABLE karena Ent tidak dukung GENERATED column)
		field.Time("att_date").
			Optional().
			Nillable().
			SchemaType(map[string]string{
				dialect.Postgres: "date",
			}),
	}
}

func (Attendance) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
