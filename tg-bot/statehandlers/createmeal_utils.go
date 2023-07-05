package statehandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
)

// MarshalBinary converts the createMealState to its binary representation.
func (cms createMealState) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(cms)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return nil, fmt.Errorf("failed to marshal CreateMealState: %v", err)
	}
	return data, nil
}

// UnmarshalBinary converts the binary representation back to createMealState.
func (cms *createMealState) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, cms)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return fmt.Errorf("failed to unmarshal CreateMealState: %v", err)
	}
	return nil
}

// SetUserState sets the create meal state for a user in Redis.
func (cms *CreateMealStateHandler) SetUserState(ctx context.Context, userID int64) error {
	state, err := cms.LocalState.MarshalBinary()
	if err != nil {
		return err
	}

	key := fmt.Sprintf(createMealStateKey, userID)

	err = cms.Client.Set(ctx, key, state, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserState deletes the create meal state for a user from Redis.
func (cms *CreateMealStateHandler) DeleteUserState(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(createMealStateKey, userID)
	err := cms.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserState gets the create meal state for a user from Redis.
func (cms *CreateMealStateHandler) GetUserState(ctx context.Context, userID int64) (createMealState, error) {
	key := fmt.Sprintf(createMealStateKey, userID)
	stateStr, err := cms.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return NoCreateMealState, nil
		}
		return NoCreateMealState, err
	}

	var state createMealState
	err = state.UnmarshalBinary([]byte(stateStr))
	if err != nil {
		return NoCreateMealState, err
	}

	return state, nil
}

// SetUserData sets the create meal data for a user in Redis.
func (cms *CreateMealStateHandler) SetUserData(ctx context.Context, userID int64) error {
	keyName := fmt.Sprintf(createMealNameKey, userID)
	err := cms.Client.Set(ctx, keyName, cms.Name, 0).Err()
	if err != nil {
		return err
	}

	keyTime := fmt.Sprintf(createMealTimeKey, userID)
	err = cms.Client.Set(ctx, keyTime, cms.Time, 0).Err()
	if err != nil {
		return err
	}

	keyRecipes := fmt.Sprintf(createMealRecipesKey, userID)
	recipesJSON, err := json.Marshal(cms.Recipes)
	if err != nil {
		return err
	}
	err = cms.Client.Set(ctx, keyRecipes, recipesJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserData deletes the create meal data for a user from Redis.
func (cms *CreateMealStateHandler) DeleteUserData(ctx context.Context, userID int64) error {
	keyName := fmt.Sprintf(createMealNameKey, userID)
	err := cms.Client.Del(ctx, keyName).Err()
	if err != nil {
		return err
	}

	keyTime := fmt.Sprintf(createMealTimeKey, userID)
	err = cms.Client.Del(ctx, keyTime).Err()
	if err != nil {
		return err
	}

	keyRecipes := fmt.Sprintf(createMealRecipesKey, userID)
	err = cms.Client.Del(ctx, keyRecipes).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserData gets the create meal data for a user from Redis.
func (cms *CreateMealStateHandler) GetUserData(ctx context.Context, userID int64) (name string, time string, recipes []int, err error) {
	keyName := fmt.Sprintf(createMealNameKey, userID)
	name, err = cms.Client.Get(ctx, keyName).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", nil, nil
		}
		return "", "", nil, err
	}

	keyTime := fmt.Sprintf(createMealTimeKey, userID)
	time, err = cms.Client.Get(ctx, keyTime).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", nil, nil
		}
		return "", "", nil, err
	}

	keyRecipes := fmt.Sprintf(createMealRecipesKey, userID)
	recipesJSON, err := cms.Client.Get(ctx, keyRecipes).Result()
	if err != nil {
		if err == redis.Nil {
			return name, time, nil, nil
		}
		return "", "", nil, err
	}

	err = json.Unmarshal([]byte(recipesJSON), &recipes)
	if err != nil {
		return "", "", nil, err
	}

	return name, time, recipes, nil
}
