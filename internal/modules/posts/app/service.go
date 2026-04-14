package app

import (
	"context"
	"time"

	"hexagonalapp/internal/modules/posts/domain"
	"hexagonalapp/internal/modules/posts/ports"
	"hexagonalapp/internal/platform/id"
)

type Service struct {
	repo   ports.Repository
	cache  ports.Cache
	mirror ports.MirrorRepository
}

func New(repo ports.Repository, cache ports.Cache, mirror ports.MirrorRepository) *Service {
	return &Service{repo: repo, cache: cache, mirror: mirror}
}

type CreatePostInput struct {
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *Service) CreatePost(ctx context.Context, in CreatePostInput) (domain.Post, error) {
	now := time.Now().UTC()
	post := domain.Post{
		ID:        id.New("pst"),
		UserID:    in.UserID,
		Title:     in.Title,
		Content:   in.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, post); err != nil {
		return domain.Post{}, err
	}
	_ = s.cache.Set(ctx, post)
	_ = s.mirror.Upsert(ctx, post)
	return post, nil
}

func (s *Service) GetPost(ctx context.Context, postID string) (domain.Post, error) {
	if post, ok, err := s.cache.Get(ctx, postID); err == nil && ok {
		return post, nil
	}
	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return domain.Post{}, err
	}
	_ = s.cache.Set(ctx, post)
	return post, nil
}

func (s *Service) ListPosts(ctx context.Context) ([]domain.Post, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListPostsByUser(ctx context.Context, userID string) ([]domain.Post, error) {
	return s.repo.ListByUserID(ctx, userID)
}
