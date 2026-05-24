package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type CronZK struct {
	ent.Schema
}

func (CronZK) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cronzk"},
	}
}

func (CronZK) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("spec"), // Cron specification (e.g., "0 * * * *")
		field.String("command"),
		field.Bool("enabled").Default(true),
		field.Time("last_run_at").Nillable().Optional(),
		field.Time("next_run_at").Nillable().Optional(),
	}
}

// Edges of the CronJob.
func (CronZK) Edges() []ent.Edge {
	return nil
}

func (CronZK) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
