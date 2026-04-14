package mongodb

import (
	"context"
	"time"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"hexagonalapp/internal/modules/user/domain"
)

type EventDocument struct {
	ID        string    `bson:"_id" json:"id"`
	Event     string    `bson:"event" json:"event"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Audit struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *Audit {
	return &Audit{collection: collection}
}

func (a *Audit) Record(ctx context.Context, event string, user domain.User) error {
	  user_ID := strconv.FormatUint(uint64(user.ID), 10)
	doc := EventDocument{
		ID:        user_ID + ":" + event,
		Event:     event,
		UserID:    user_ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Now().UTC(),
	}
	_, err := a.collection.InsertOne(ctx, doc)
	return err
}
