package app_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	apppkg "hexagonalapp/internal/modules/posts/app"
	"hexagonalapp/internal/modules/posts/domain"
)

// --- Fake implementations for ports interfaces ---

type fakeRepo struct {
	storage   map[string]domain.Post
	nextID    uint
	getCalled int
	createErr error
	updateErr error
	deleteErr error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		storage: make(map[string]domain.Post),
		nextID:  1,
	}
}

func (r *fakeRepo) Create(ctx context.Context, post *domain.Post) error {
	if r.createErr != nil {
		return r.createErr
	}
	if post.ID == 0 {
		post.ID = r.nextID
		r.nextID++
	}
	id := strconv.FormatUint(uint64(post.ID), 10)
	r.storage[id] = *post
	return nil
}

func (r *fakeRepo) Update(ctx context.Context, post domain.Post) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	id := strconv.FormatUint(uint64(post.ID), 10)
	if _, ok := r.storage[id]; !ok {
		return errors.New("not found")
	}
	r.storage[id] = post
	return nil
}

func (r *fakeRepo) Delete(ctx context.Context, id string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	if _, ok := r.storage[id]; !ok {
		return errors.New("not found")
	}
	delete(r.storage, id)
	return nil
}

func (r *fakeRepo) GetByID(ctx context.Context, id string) (domain.Post, error) {
	r.getCalled++
	if p, ok := r.storage[id]; ok {
		return p, nil
	}
	return domain.Post{}, errors.New("not found")
}

func (r *fakeRepo) List(ctx context.Context) ([]domain.Post, error) {
	var out []domain.Post
	for _, p := range r.storage {
		out = append(out, p)
	}
	return out, nil
}

func (r *fakeRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Post, error) {
	var out []domain.Post
	for _, p := range r.storage {
		if p.UserID == userID {
			out = append(out, p)
		}
	}
	return out, nil
}

// --- Fake cache implementation ---

type fakeCache struct {
	store      map[string]domain.Post
	setCalled  int
	delCalled  int
	getReturn  domain.Post
	getOk      bool
	getErr     error
	lastSetKey string
}

func newFakeCache() *fakeCache {
	return &fakeCache{store: make(map[string]domain.Post)}
}

func (c *fakeCache) Get(ctx context.Context, id string) (domain.Post, bool, error) {
	if c.getErr != nil {
		return domain.Post{}, false, c.getErr
	}
	if p, ok := c.store[id]; ok {
		return p, true, nil
	}
	// support configured return (for test convenience)
	if c.getOk {
		return c.getReturn, true, nil
	}
	return domain.Post{}, false, nil
}

func (c *fakeCache) Set(ctx context.Context, post domain.Post) error {
	c.setCalled++
	id := strconv.FormatUint(uint64(post.ID), 10)
	c.store[id] = post
	c.lastSetKey = id
	return nil
}

func (c *fakeCache) Delete(ctx context.Context, id string) error {
	c.delCalled++
	delete(c.store, id)
	return nil
}

// --- Fake mirror implementation ---

type fakeMirror struct {
	upsertCalled int
	lastUpsert   domain.Post
	err          error
}

func (m *fakeMirror) Upsert(ctx context.Context, post domain.Post) error {
	m.upsertCalled++
	m.lastUpsert = post
	return m.err
}

// --- Tests ---

func TestCreatePost_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}

	svc := apppkg.New(repo, cache, mirror)

	input := domain.CreatePostInput{
		Title:   "Hello",
		Content: "World content",
		UserID:  "42",
	}

	post, _, err := svc.CreatePost(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if post.ID == 0 {
		t.Fatalf("expected post ID to be set, got 0")
	}

	// repo should contain the new post
	stored, err := repo.GetByID(ctx, strconv.FormatUint(uint64(post.ID), 10))
	if err != nil {
		t.Fatalf("repo missing created post: %v", err)
	}
	if stored.Title != input.Title || stored.Content != input.Content {
		t.Fatalf("stored post mismatch: want %v/%v got %v/%v", input.Title, input.Content, stored.Title, stored.Content)
	}

	// cache should have been set
	if cache.setCalled == 0 {
		t.Fatalf("expected cache.Set to be called")
	}
	// mirror upsert should be called
	if mirror.upsertCalled == 0 {
		t.Fatalf("expected mirror.Upsert to be called")
	}
}

func TestCreatePost_ValidationFailure(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}
	svc := apppkg.New(repo, cache, mirror)

	// invalid: title too short (validator requires min=2)
	input := domain.CreatePostInput{
		Title:   "A",
		Content: "ok",
		UserID:  "1",
	}

	_, _, err := svc.CreatePost(ctx, input)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}

func TestGetPost_CacheHit(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}

	// prepare cache to return a post
	cached := domain.Post{
		ID:      7,
		Title:   "cached",
		Content: "c",
		UserID:  "u1",
	}
	cache.store[strconv.FormatUint(uint64(cached.ID), 10)] = cached

	svc := apppkg.New(repo, cache, mirror)

	got, err := svc.GetPost(ctx, strconv.FormatUint(uint64(cached.ID), 10))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != cached.Title {
		t.Fatalf("expected cached title %q, got %q", cached.Title, got.Title)
	}
	// repo should not be hit
	if repo.getCalled != 0 {
		t.Fatalf("expected repo.GetByID not called, got %d", repo.getCalled)
	}
}

func TestGetPost_CacheMiss(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}

	// prepare repo with a post
	p := domain.Post{
		ID:        9,
		Title:     "fromrepo",
		Content:   "r",
		UserID:    "u2",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_ = repo.Create(ctx, &p)

	svc := apppkg.New(repo, cache, mirror)

	got, err := svc.GetPost(ctx, strconv.FormatUint(uint64(p.ID), 10))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != p.Title {
		t.Fatalf("expected title %q got %q", p.Title, got.Title)
	}
	// cache should now contain it
	if cache.setCalled == 0 {
		t.Fatalf("expected cache.Set to be called on cache miss")
	}
}

func TestUpdatePost_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}

	// create an initial post
	initial := domain.Post{
		ID:        11,
		UserID:    "user11",
		Title:     "old",
		Content:   "oldcontent",
		CreatedAt: time.Now().UTC().Add(-time.Hour),
		UpdatedAt: time.Now().UTC().Add(-time.Hour),
	}
	_ = repo.Create(ctx, &initial)

	svc := apppkg.New(repo, cache, mirror)

	input := domain.CreatePostInput{
		UserID:  "user11",
		Title:   "new title",
		Content: "new content",
	}

	updated, _, err := svc.UpdatePost(ctx, strconv.FormatUint(uint64(initial.ID), 10), input)
	if err != nil {
		t.Fatalf("expected no error updating post, got %v", err)
	}
	if updated.Title != input.Title || updated.Content != input.Content {
		t.Fatalf("update did not apply properly")
	}
	// cache and mirror should be updated
	if cache.setCalled == 0 {
		t.Fatalf("expected cache.Set on update")
	}
	if mirror.upsertCalled == 0 {
		t.Fatalf("expected mirror.Upsert on update")
	}
}

func TestDeletePost_Success(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}

	// create entry and seed cache
	p := domain.Post{
		ID:      13,
		UserID:  "u13",
		Title:   "t13",
		Content: "c13",
	}
	_ = repo.Create(ctx, &p)
	cache.Set(ctx, p)

	svc := apppkg.New(repo, cache, mirror)

	err := svc.DeletePost(ctx, strconv.FormatUint(uint64(p.ID), 10))
	if err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
	// cache.Delete should be called
	if cache.delCalled == 0 {
		t.Fatalf("expected cache.Delete to be called on delete")
	}
	// repo should no longer contain it
	_, err = repo.GetByID(ctx, strconv.FormatUint(uint64(p.ID), 10))
	if err == nil {
		t.Fatalf("expected repo.GetByID to fail after delete")
	}
}

func TestListPosts(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRepo()
	cache := newFakeCache()
	mirror := &fakeMirror{}

	// seed repo with entries
	for i := 1; i <= 3; i++ {
		p := domain.Post{
			ID:      uint(i + 20),
			UserID:  "u",
			Title:   "t",
			Content: "c",
		}
		_ = repo.Create(ctx, &p)
	}

	svc := apppkg.New(repo, cache, mirror)
	list, err := svc.ListPosts(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(list))
	}
}
