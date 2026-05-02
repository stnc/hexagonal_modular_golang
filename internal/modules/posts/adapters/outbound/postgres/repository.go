package postgres

import (
	"context"
	"errors"
	"time"

	"hexagonalapp/internal/modules/posts/domain"

	"gorm.io/gorm"
)

type PostModel struct {
	ID        uint      `gorm:"primaryKey;auto_increment;column:id"`
	UserID    string    `gorm:"column:user_id;index"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PostModel) TableName() string { return "posts" }

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&PostModel{})
}

func (r *Repository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Table("posts").Create(post).Error
}

func (r *Repository) Update(ctx context.Context, post domain.Post) error {
	m := toModel(post)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&PostModel{}, "id = ?", id).Error
}

func (r *Repository) Upsert(ctx context.Context, post domain.Post) error {
	m := toModel(post)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *Repository) GetByID(ctx context.Context, id string) (domain.Post, error) {
	var m PostModel
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Post{}, err
		}
		return domain.Post{}, err
	}
	return toDomain(m), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Post, error) {
	var rows []PostModel
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Post, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomain(row))
	}
	return items, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]domain.Post, error) {
	var rows []PostModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Post, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomain(row))
	}
	return items, nil
}

func toModel(p domain.Post) PostModel {
	return PostModel{ID: p.ID, UserID: p.UserID, Title: p.Title, Content: p.Content, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func toDomain(m PostModel) domain.Post {
	return domain.Post{ID: m.ID, UserID: m.UserID, Title: m.Title, Content: m.Content, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
