package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Department struct {
	ent.Schema
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "department"},
	}
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("depart_code").Optional(),
		field.String("depart_name"),
		field.String("status").Default("0"),
		field.String("ket").Optional(),
	}
}

func (Department) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
