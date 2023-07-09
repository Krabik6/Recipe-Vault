package manager

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type JwtManager struct {
	Client     *redis.Client
	Expiration time.Duration //expiration
}

func NewJwtManager(client *redis.Client, expiration time.Duration) *JwtManager {
	return &JwtManager{Client: client, Expiration: expiration}
}

// SetUserJWTToken sets the JWT token for a user in Redis.
func (jm *JwtManager) SetUserJWTToken(ctx context.Context, userID int64, token string) error {
	key := fmt.Sprintf(model.UserJWTTokenKey, userID)

	exp := jm.Expiration
	log.Println(exp)
	err := jm.Client.Set(ctx, key, token, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserJWTToken gets the JWT token for a user from Redis.
func (jm *JwtManager) GetUserJWTToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf(model.UserJWTTokenKey, userID)
	token, err := jm.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return token, nil
}

// DeleteUserJWTToken deletes the JWT token for a user from Redis.
func (jm *JwtManager) DeleteUserJWTToken(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(model.UserJWTTokenKey, userID)
	err := jm.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// CheckLoggedIn func to check if user is logged in (returns true if logged in)
func (jm *JwtManager) CheckLoggedIn(ctx context.Context, userID int64) (bool, error) {
	// Getting user state from manager
	loggedIn, err := jm.GetUserJWTToken(ctx, userID)
	if err != nil {
		return false, err
	}
	if loggedIn != "" {
		return true, nil
	}
	return false, nil
}
