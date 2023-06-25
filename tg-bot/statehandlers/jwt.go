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
