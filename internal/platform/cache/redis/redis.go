package redis

import redisclient "github.com/redis/go-redis/v9"

func New(addr, password string, db int) *redisclient.Client {
	return redisclient.NewClient(&redisclient.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
