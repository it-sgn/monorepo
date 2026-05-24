package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"entgo.io/ent/schema/mixin"
)

type TimeMixin struct {
	mixin.Schema
}

func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("created_by").Optional(),
		field.String("updated_by").Optional(),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Optional(),
	}
}
