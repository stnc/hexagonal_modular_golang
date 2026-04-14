package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"hexagonalapp/internal/modules/posts/domain"
)

type PostDocument struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	Title     string    `bson:"title" json:"title"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Repository struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

func (r *Repository) Create(ctx context.Context, post domain.Post) error {
	_, err := r.collection.InsertOne(ctx, toDocument(post))
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (domain.Post, error) {
	var doc PostDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Post{}, err
		}
		return domain.Post{}, err
	}
	return toDomain(doc), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []PostDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	posts := make([]domain.Post, 0, len(docs))
	for _, doc := range docs {
		posts = append(posts, toDomain(doc))
	}
	return posts, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]domain.Post, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []PostDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	posts := make([]domain.Post, 0, len(docs))
	for _, doc := range docs {
		posts = append(posts, toDomain(doc))
	}
	return posts, nil
}

func toDocument(post domain.Post) PostDocument {
	return PostDocument{ID: post.ID, UserID: post.UserID, Title: post.Title, Content: post.Content, CreatedAt: post.CreatedAt, UpdatedAt: post.UpdatedAt}
}

func toDomain(doc PostDocument) domain.Post {
	return domain.Post{ID: doc.ID, UserID: doc.UserID, Title: doc.Title, Content: doc.Content, CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt}
}
