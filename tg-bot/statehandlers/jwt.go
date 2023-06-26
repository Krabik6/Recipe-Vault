package statehandlers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const userJWTTokenKey = "user_jwt_token:%d"

// SetUserJWTToken sets the JWT token for a user in Redis.
func (sh *StateHandler) SetUserJWTToken(ctx context.Context, userID int64, token string) error {
	key := fmt.Sprintf(userJWTTokenKey, userID)
	// expiration - 7 days type time.Duration
	exp := time.Duration(7 * 24 * time.Hour)

	err := sh.Client.Set(ctx, key, token, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserJWTToken gets the JWT token for a user from Redis.
func (sh *StateHandler) GetUserJWTToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf(userJWTTokenKey, userID)
	token, err := sh.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return token, nil
}

// DeleteUserJWTToken deletes the JWT token for a user from Redis.
func (sh *StateHandler) DeleteUserJWTToken(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(userJWTTokenKey, userID)
	err := sh.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// CheckLoggedIn func to check if user is logged in (returns true if logged in)
func (sh *StateHandler) CheckLoggedIn(ctx context.Context, userID int64) (bool, error) {
	// Getting user state from redis
	loggedIn, err := sh.GetUserJWTToken(ctx, userID)
	if err != nil {
		return false, err
	}
	if loggedIn != "" {
		return true, nil
	}
	return false, nil
}
