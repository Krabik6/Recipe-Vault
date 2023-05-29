package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

// CRD is the interface for a Redis cache client.
type CRD interface {
	SetKey(key string, value interface{}, expiration int64) error    // Set a key-value pair with expiration.
	GetKey(key string) (string, error)                               // Get the value for a key.
	DeleteKey(key string) error                                      // Delete a key.
	SetObject(key string, value interface{}, expiration int64) error // Set an object in the cache with expiration.
	GetObject(key string, value interface{}) error                   // Get an object from the cache.

}

// CRDRedis is the implementation of the CRD interface using go-redis/v9 Redis client.
type CRDRedis struct {
	redisClient *redis.Client
}

// NewCRDRedis creates a new Redis cache client instance.
func NewCRDRedis(client redis.Client) *CRDRedis {
	return &CRDRedis{redisClient: &client}
}

// SetObject sets an object in the Redis cache with expiration.
func (crd *CRDRedis) SetObject(key string, value interface{}, expiration int64) error {
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = crd.redisClient.Set(ctx, key, jsonValue, time.Duration(expiration)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetObject gets the value of an object from the Redis cache.
func (crd *CRDRedis) GetObject(key string, value interface{}) error {
	ctx := context.Background()
	jsonValue, err := crd.redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonValue), value)
	if err != nil {
		return err
	}
	return nil
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
func (crd *CRDRedis) DeleteKey(key string) error {
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
