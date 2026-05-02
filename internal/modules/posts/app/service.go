package app

import (
	"context"
	"time"

	"hexagonalapp/internal/modules/posts/domain"
	"hexagonalapp/internal/modules/posts/ports"

	"github.com/gofiber/fiber/v3"
	// "hexagonalapp/internal/platform/id"
)

type Service struct {
	repo   ports.Repository
	cache  ports.Cache
	mirror ports.MirrorRepository
}

func New(repo ports.Repository, cache ports.Cache, mirror ports.MirrorRepository) *Service {
	return &Service{repo: repo, cache: cache, mirror: mirror}
}

func (s *Service) CreatePost(ctx context.Context, input domain.CreatePostInput) (*domain.Post, fiber.Map, error) {
	input = domain.NormalizeInput(input)
	if err := domain.ValidateInput(input); err != nil {
		return &domain.Post{}, fiber.Map{}, err
	}

	now := time.Now().UTC()
	post := &domain.Post{
		ID:        input.ID,
		UserID:    input.UserID,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, post); err != nil {
		return nil, fiber.Map{}, err
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, *post)
	}
	if s.mirror != nil {
		_ = s.mirror.Upsert(ctx, *post)
	}
	return post, fiber.Map{}, nil
}

func (s *Service) GetPost(ctx context.Context, postID string) (domain.Post, error) {
	if s.cache != nil {
		if post, ok, err := s.cache.Get(ctx, postID); err == nil && ok {
			return post, nil
		}
	}

	post, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return domain.Post{}, err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, post)
	}
	return post, nil
}

func (s *Service) ListPosts(ctx context.Context) ([]domain.Post, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListPostsByUser(ctx context.Context, userID string) ([]domain.Post, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *Service) UpdatePost(ctx context.Context, postID string, input domain.CreatePostInput) (*domain.Post, fiber.Map, error) {
	input = domain.NormalizeInput(input)
	if err := domain.ValidateInput(input); err != nil {
		return &domain.Post{}, fiber.Map{}, err
	}

	current, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return nil, fiber.Map{}, err
	}

	current.UserID = input.UserID
	current.Title = input.Title
	current.Content = input.Content
	current.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, current); err != nil {
		return nil, fiber.Map{}, err
	}
	if s.cache != nil {
		_ = s.cache.Set(ctx, current)
	}
	if s.mirror != nil {
		_ = s.mirror.Upsert(ctx, current)
	}

	return &current, fiber.Map{}, nil
}

func (s *Service) DeletePost(ctx context.Context, postID string) error {
	if err := s.repo.Delete(ctx, postID); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.Delete(ctx, postID)
	}
	return nil
}
