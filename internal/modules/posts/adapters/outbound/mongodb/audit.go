package mongodb

import (
	"context"
	"time"

	"hexagonalapp/internal/modules/posts/domain"
	conventorLib "hexagonalapp/internal/platform/helpers/stnccollection"

	"go.mongodb.org/mongo-driver/v2/mongo"
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
	ID := conventorLib.UintToString(post.ID)
	_, err := a.collection.InsertOne(ctx, EventDocument{
		ID:        ID + ":" + event,
		Event:     event,
		PostID:    ID,
		UserID:    post.UserID,
		Title:     post.Title,
		CreatedAt: time.Now().UTC(),
	})
	return err
}
