package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Assignment defines the penempatan karyawan ke jabatan
type Assignment struct {
	ent.Schema
}

func (Assignment) Fields() []ent.Field {
	return []ent.Field{
		field.String("employee_id").
			NotEmpty().
			Comment("kode karyawan"),

		field.Int64("position_id").
			// NotEmpty().
			Comment("kode jabatan"),
		field.Time("start_date").
			SchemaType(map[string]string{
				"postgres": "timestamptz",
			}).
			Default(time.Now),
		// field.Time("start_date").
		// 	Default(time.Now),

		field.Time("end_date").
			SchemaType(map[string]string{
				"postgres": "timestamptz",
			}).
			Optional().
			Nillable(),
	}
}

func (Assignment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Assignment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("position", Position.Type).
			Ref("assignments").
			Required().
			Field("position_id").
			Unique(),
	}
}
