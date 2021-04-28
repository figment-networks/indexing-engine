package datalake

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisStorage struct {
	client     *redis.Client
	expiration time.Duration
}

var _ Storage = (*redisStorage)(nil)

// NewRedisStorage creates a Redis storage
func NewRedisStorage(addr string, exp time.Duration) Storage {
	opts := redis.Options{Addr: addr}

	return &redisStorage{
		client:     redis.NewClient(&opts),
		expiration: exp,
	}
}

func (rs *redisStorage) Store(data []byte, path ...string) error {
	ctx := context.Background()
	key := rs.storageKey(path)

	return rs.client.Set(ctx, key, data, rs.expiration).Err()
}

func (rs *redisStorage) IsStored(path ...string) (bool, error) {
	ctx := context.Background()
	key := rs.storageKey(path)

	cmd := rs.client.Exists(ctx, key)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() == 1, nil
}

func (rs *redisStorage) Retrieve(path ...string) ([]byte, error) {
	ctx := context.Background()
	key := rs.storageKey(path)

	return rs.client.Get(ctx, key).Bytes()
}

func (rs *redisStorage) storageKey(path []string) string {
	return strings.Join(path, ":")
}
