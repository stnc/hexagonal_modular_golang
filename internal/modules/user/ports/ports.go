package ports

import (
	"context"
	"github.com/gofiber/fiber/v3"

	"hexagonalapp/internal/modules/user/domain"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	ExistsByEmail(ctx context.Context, email string, excludeID uint) (bool, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	ListDataTable(c fiber.Ctx, ctx context.Context) ([]domain.User, int64, int64, error)
	ListPagination(ctx context.Context, postsPerPage, offset int) ([]domain.User, error)
	Count(*int64)
}

type Cache interface {
	Get(ctx context.Context, id string) (domain.User, bool, error)
	Set(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
}

type Audit interface {
	Record(ctx context.Context, event string, user domain.User) error
}
