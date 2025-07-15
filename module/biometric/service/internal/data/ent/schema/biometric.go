package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Biometric struct {
	ent.Schema
}

func (Biometric) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "biometric"},
	}
}

func (Biometric) Fields() []ent.Field {
	return []ent.Field{
		field.String("fingercode").Optional(),

		field.String("finger0").Optional().
			Optional(). // boleh kosong (nullable)
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger1").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger2").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger3").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger4").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger5").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger6").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger7").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger8").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
		field.String("finger9").Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "text", // paksa gunakan TEXT (meskipun defaultnya juga TEXT jika tanpa MaxLen)
			}),
	}
}

func (Biometric) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
