package ports

import (
	"context"

	"hexagonalapp/internal/modules/posts/domain"
)

type Repository interface {
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post domain.Post) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (domain.Post, error)
	List(ctx context.Context) ([]domain.Post, error)
	ListByUserID(ctx context.Context, userID string) ([]domain.Post, error)
}

type Cache interface {
	Get(ctx context.Context, id string) (domain.Post, bool, error)
	Set(ctx context.Context, post domain.Post) error
	Delete(ctx context.Context, id string) error
}

type MirrorRepository interface {
	Upsert(ctx context.Context, post domain.Post) error
}
