package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"hexagonalapp/internal/modules/posts/domain"

	redisclient "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redisclient.Client
	ttl    time.Duration
}

func New(client *redisclient.Client) *Cache {
	return &Cache{client: client, ttl: 5 * time.Minute}
}

func (c *Cache) key(id string) string { return fmt.Sprintf("post:%s", id) }

func (c *Cache) Get(ctx context.Context, id string) (domain.Post, bool, error) {
	data, err := c.client.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		if err == redisclient.Nil {
			return domain.Post{}, false, nil
		}
		return domain.Post{}, false, err
	}
	var post domain.Post
	if err := json.Unmarshal(data, &post); err != nil {
		return domain.Post{}, false, err
	}
	return post, true, nil
}

func (c *Cache) Set(ctx context.Context, post domain.Post) error {
	ID := strconv.FormatUint(uint64(post.ID), 10)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(ID), data, c.ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
