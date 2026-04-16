package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	redisclient "github.com/redis/go-redis/v9"
	"hexagonalapp/internal/modules/user/domain"
)

type Cache struct {
	client *redisclient.Client
	ttl    time.Duration
}

func New(client *redisclient.Client) *Cache {
	return &Cache{client: client, ttl: 5 * time.Minute}
}

func (c *Cache) key(id string) string { return fmt.Sprintf("user:%s", id) }

func (c *Cache) Get(ctx context.Context, id string) (domain.User, bool, error) {
	data, err := c.client.Get(ctx, c.key(id)).Bytes()
	// fmt.Println("stoop")
	// fmt.Println(c.key(id))
	// fmt.Println(data)
	// // return domain.User{}, true, nil
	if err != nil {
		if err == redisclient.Nil {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}
	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return domain.User{}, false, err
	}
	
	// fmt.Println(user)

	return user, true, nil
}

func (c *Cache) Set(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	// fmt.Println(data)
	if err != nil {
		return err
	}
	user_ID := strconv.FormatUint(uint64(user.ID), 10)
	// fmt.Println(user_ID)
	return c.client.Set(ctx, c.key(user_ID), data, c.ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
