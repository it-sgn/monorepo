package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Position defines the Jabatan (role-related) structure
type Perusahaan struct {
	ent.Schema
}

func (Perusahaan) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Comment("position_id"),
		field.String("kode_perusahaan").Unique().Comment("kode_perusahaan"),
		field.String("nama_perusahaan").
			NotEmpty().
			Comment("Nama perusahaan, e.g., PT. SGN"),
		field.String("kode_cabang").
			NotEmpty().Unique().
			Comment("kode cabang contoh SG33"),

		field.String("cabang").
			NotEmpty().Unique().
			Comment("Nama cabang"),

		field.String("alamat").
			Optional().
			Comment("Alamat"),

		field.String("telp").
			Optional().
			Nillable().
			Comment("telp"),
		field.String("email").
			Optional().
			Nillable().
			Comment("email"),
		field.String("logo").Nillable().Optional(),
	}
}

func (Perusahaan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
