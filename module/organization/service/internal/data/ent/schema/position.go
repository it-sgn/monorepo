package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Position defines the Jabatan (role-related) structure
type Position struct {
	ent.Schema
}

func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("position_id"),
		field.String("position_code").Unique().Comment("position_code"),
		field.String("name").
			NotEmpty().
			Comment("Nama jabatan, e.g., Admin HRD"),

		field.String("role_name").
			NotEmpty().
			Comment("Role name di Keycloak"),

		field.String("department_code").
			Optional().
			Comment("Relasi ke Department"),

		field.String("reports_to_position_id").
			Optional().
			Nillable().
			Comment("Jabatan atasan"),
	}
}

func (Position) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assignments", Assignment.Type),
	}
}

func (Position) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
