package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Employers struct {
	ent.Schema
}

func (Employers) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "employers"},
	}
}

func (Employers) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("nosap").Comment("").Nillable(),
		field.String("nip").Comment("").Nillable(),
		field.String("karyacode").Comment("").Nillable(),
		field.String("karyaname").Comment(""),
		field.String("disp_name").Comment(""),
		field.String("pass_mesin").Comment("").Nillable(),
		field.String("rfid_card").Comment("").Nillable(),
		field.String("kode_finger").Comment(""),
		field.String("depart_code").Comment(""),
		field.Int32("status").Default(0),
	}
}

func (Employers) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
