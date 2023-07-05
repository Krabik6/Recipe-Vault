package manager

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type UserStateManager struct {
	Client *redis.Client
}

// NewRedisUserStateManager - constructor for UserStateManager
func NewRedisUserStateManager(client *redis.Client) *UserStateManager {
	return &UserStateManager{
		Client: client,
	}
}

func (r *UserStateManager) SetUserState(ctx context.Context, userID int64, state model.State) error {
	// Forming the key for the user
	key := fmt.Sprintf(model.UserState, userID)

	// Setting the state value in Redis
	err := r.Client.Set(ctx, key, state, 0).Err()
	if err != nil {
		return err // Return the error if there was a problem setting the state in Redis
	}

	return nil // Return nil to indicate success
}

func (r *UserStateManager) GetUserState(ctx context.Context, userID int64) (model.State, error) {
	// Forming the key for the user
	key := fmt.Sprintf(model.UserState, userID)

	// Getting the state value from Redis
	stateStr, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		// Handling the case when the value is not found in Redis
		if err == redis.Nil {
			return model.NoState, nil // Return NoState if the state is not found
		}
		return model.NoState, err // Return an error if there is any other Redis error
	}

	state, err := strconv.Atoi(stateStr)
	if err != nil {
		return model.NoState, err // Return an error if there is an error in converting the state value
	}

	return model.State(state), nil // Return the user state
}

func (r *UserStateManager) DeleteUserState(ctx context.Context, userID int64) error {
	// Forming the key for the user
	key := fmt.Sprintf(model.UserState, userID)

	// Deleting the state value from Redis
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return err // Return the error if there was a problem deleting the state from Redis
	}

	return nil // Return nil to indicate success
}
