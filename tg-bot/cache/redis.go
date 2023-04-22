package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// CRD is the interface for a Redis cache client.
type CRD interface {
	SetKey(key string, value interface{}, expiration int64) error // Set a key-value pair with expiration.
	GetKey(key string) (string, error)                            // Get the value for a key.
	DelKey(key string) error                                      // Delete a key.
}

// CRDRedis is the implementation of the CRD interface using go-redis/v9 Redis client.
type CRDRedis struct {
	redisClient *redis.Client
}

// NewCRDRedis creates a new Redis cache client instance.
func NewCRDRedis(client redis.Client) *CRDRedis {
	return &CRDRedis{redisClient: &client}
}

// SetKey sets a key-value pair in the Redis cache with expiration.
func (crd *CRDRedis) SetKey(key string, value interface{}, expiration int64) error {
	ctx := context.Background()
	err := crd.redisClient.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetKey gets the value for a key in the Redis cache.
func (crd *CRDRedis) GetKey(key string) (string, error) {
	ctx := context.Background()
	val, err := crd.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// DelKey deletes a key from the Redis cache.
func (crd *CRDRedis) DelKey(key string) error {
	ctx := context.Background()
	err := crd.redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// Repository is a wrapper for the CRD interface implementation.
type Repository struct {
	CRD
}

// NewCacheRepository creates a new cache repository using a Redis cache client.
func NewCacheRepository(client *redis.Client) *Repository {
	return &Repository{CRD: NewCRDRedis(*client)}
}
