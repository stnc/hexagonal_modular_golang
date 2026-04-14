package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"hexagonalapp/internal/modules/posts/domain"
)

type Audit struct {
	collection *mongo.Collection
}

func NewAudit(collection *mongo.Collection) *Audit {
	return &Audit{collection: collection}
}

type EventDocument struct {
	ID        string    `bson:"_id"`
	Event     string    `bson:"event"`
	PostID    string    `bson:"post_id"`
	UserID    string    `bson:"user_id"`
	Title     string    `bson:"title"`
	CreatedAt time.Time `bson:"created_at"`
}

func (a *Audit) Record(ctx context.Context, event string, post domain.Post) error {
	_, err := a.collection.InsertOne(ctx, EventDocument{
		ID:        post.ID + ":" + event,
		Event:     event,
		PostID:    post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		CreatedAt: time.Now().UTC(),
	})
	return err
}
