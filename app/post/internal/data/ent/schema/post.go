package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

func (Post) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sns_post"},
		// entsql.WithComments(true),
	}
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").SchemaType(map[string]string{
			dialect.SQLite: "VARCHAR(36) PRIMARY KEY",
		}).Unique().DefaultFunc(func() string {
			return uuid.New().String()
		}).Immutable().NotEmpty().MaxLen(36).Comment("文章ID"),
		field.String("title").SchemaType(map[string]string{
			dialect.SQLite: "VARCHAR(255)",
		}).MinLen(1).MaxLen(255).Comment("标题"),
		field.String("content").Default("").Comment("正文"),
		field.Time("created_at").SchemaType(map[string]string{
			dialect.SQLite: "DATETIME",
		}).Default(time.Now).Immutable().Annotations(entsql.Default("CURRENT_TIMESTAMP")).Comment("创建时间"),
		field.Time("updated_at").SchemaType(map[string]string{
			dialect.SQLite: "DATETIME",
		}).Default(time.Now).UpdateDefault(time.Now).Annotations(entsql.Default("CURRENT_TIMESTAMP")).Comment("更新时间"),
		field.Int("deleted_at").Default(0).Comment("删除时间"),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return nil
}

func (Post) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title", "deleted_at").Unique(),
	}
}
