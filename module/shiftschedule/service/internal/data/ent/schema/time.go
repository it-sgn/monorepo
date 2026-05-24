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
		field.Time("create_time").
			Immutable().
			Default(time.Now),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("delete_time").
			Optional(),
	}
}
