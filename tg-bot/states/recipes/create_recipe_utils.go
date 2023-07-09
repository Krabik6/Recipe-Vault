package recipes

// file name: createrecipestatehandler_utils.go
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
)

const (
	createRecipeStateKey         = "create_recipe_state:%d"
	createRecipeTitleKey         = "create_recipe_title:%d"
	createRecipeDescriptionKey   = "create_recipe_description:%d"
	createRecipeIsPublicKey      = "create_recipe_isPublic:%d"
	createRecipeCostKey          = "create_recipe_cost:%d"
	createRecipeTimeToPrepareKey = "create_recipe_timeToPrepare:%d"
	createRecipeHealthyKey       = "create_recipe_healthy:%d"
)

func (crs createRecipeState) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(crs)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return nil, fmt.Errorf("failed to marshal CreateRecipeState: %v", err)
	}
	return data, nil
}

func (crs *createRecipeState) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, crs)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return fmt.Errorf("failed to unmarshal CreateRecipeState: %v", err)
	}
	return nil
}

// SetUserState sets the create recipe state for a user in Redis.
func (crs *CreateRecipeService) SetUserState(ctx context.Context, userID int64) error {
	state, err := crs.LocalState.MarshalBinary()
	if err != nil {
		return err
	}

	key := fmt.Sprintf(createRecipeStateKey, userID)

	err = crs.Client.Set(ctx, key, state, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserState deletes the create recipe state for a user from Redis.
func (crs *CreateRecipeService) DeleteUserState(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(createRecipeStateKey, userID)
	err := crs.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserState gets the create recipe state for a user from Redis.
func (crs *CreateRecipeService) GetUserState(ctx context.Context, userID int64) (createRecipeState, error) {
	key := fmt.Sprintf(createRecipeStateKey, userID)
	stateStr, err := crs.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return NoCreateRecipeState, nil
		}
		return NoCreateRecipeState, err
	}

	var state createRecipeState
	err = state.UnmarshalBinary([]byte(stateStr))
	if err != nil {
		return NoCreateRecipeState, err
	}

	return state, nil
}

// SetUserData sets the create recipe data for a user in Redis.
func (crs *CreateRecipeService) SetUserData(ctx context.Context, userID int64) error {
	keyTitle := fmt.Sprintf(createRecipeTitleKey, userID)
	err := crs.Client.Set(ctx, keyTitle, crs.Title, 0).Err()
	if err != nil {
		return err
	}

	keyDescription := fmt.Sprintf(createRecipeDescriptionKey, userID)
	err = crs.Client.Set(ctx, keyDescription, crs.Description, 0).Err()
	if err != nil {
		return err
	}

	keyIsPublic := fmt.Sprintf(createRecipeIsPublicKey, userID)
	err = crs.Client.Set(ctx, keyIsPublic, crs.IsPublic, 0).Err()
	if err != nil {
		return err
	}

	keyCost := fmt.Sprintf(createRecipeCostKey, userID)
	err = crs.Client.Set(ctx, keyCost, crs.Cost, 0).Err()
	if err != nil {
		return err
	}

	keyTimeToPrepare := fmt.Sprintf(createRecipeTimeToPrepareKey, userID)
	err = crs.Client.Set(ctx, keyTimeToPrepare, crs.TimeToPrepare, 0).Err()
	if err != nil {
		return err
	}

	keyHealthy := fmt.Sprintf(createRecipeHealthyKey, userID)
	err = crs.Client.Set(ctx, keyHealthy, crs.Healthy, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserData deletes the create recipe data for a user from Redis.
func (crs *CreateRecipeService) DeleteUserData(ctx context.Context, userID int64) error {
	keyTitle := fmt.Sprintf(createRecipeTitleKey, userID)
	err := crs.Client.Del(ctx, keyTitle).Err()
	if err != nil {
		return err
	}

	keyDescription := fmt.Sprintf(createRecipeDescriptionKey, userID)
	err = crs.Client.Del(ctx, keyDescription).Err()
	if err != nil {
		return err
	}

	keyIsPublic := fmt.Sprintf(createRecipeIsPublicKey, userID)
	err = crs.Client.Del(ctx, keyIsPublic).Err()
	if err != nil {
		return err
	}

	keyCost := fmt.Sprintf(createRecipeCostKey, userID)
	err = crs.Client.Del(ctx, keyCost).Err()
	if err != nil {
		return err
	}

	keyTimeToPrepare := fmt.Sprintf(createRecipeTimeToPrepareKey, userID)
	err = crs.Client.Del(ctx, keyTimeToPrepare).Err()
	if err != nil {
		return err
	}

	keyHealthy := fmt.Sprintf(createRecipeHealthyKey, userID)
	err = crs.Client.Del(ctx, keyHealthy).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserData gets the create recipe data for a user from Redis.
func (crs *CreateRecipeService) GetUserData(ctx context.Context, userID int64) (title string, description string, isPublic bool, cost float64, timeToPrepare int64, healthy int, err error) {
	keyTitle := fmt.Sprintf(createRecipeTitleKey, userID)
	title, err = crs.Client.Get(ctx, keyTitle).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", false, 0, 0, 0, nil
		}
		return "", "", false, 0, 0, 0, err
	}

	keyDescription := fmt.Sprintf(createRecipeDescriptionKey, userID)
	description, err = crs.Client.Get(ctx, keyDescription).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", false, 0, 0, 0, nil
		}
		return "", "", false, 0, 0, 0, err
	}

	keyIsPublic := fmt.Sprintf(createRecipeIsPublicKey, userID)
	isPublic, err = crs.Client.Get(ctx, keyIsPublic).Bool()
	if err != nil {
		return "", "", false, 0, 0, 0, err
	}

	keyCost := fmt.Sprintf(createRecipeCostKey, userID)
	cost, err = crs.Client.Get(ctx, keyCost).Float64()
	if err != nil {
		return "", "", false, 0, 0, 0, err
	}

	keyTimeToPrepare := fmt.Sprintf(createRecipeTimeToPrepareKey, userID)
	timeToPrepare, err = crs.Client.Get(ctx, keyTimeToPrepare).Int64()
	if err != nil {
		return "", "", false, 0, 0, 0, err
	}

	keyHealthy := fmt.Sprintf(createRecipeHealthyKey, userID)
	healthy, err = crs.Client.Get(ctx, keyHealthy).Int()
	if err != nil {
		return "", "", false, 0, 0, 0, err
	}

	return title, description, isPublic, cost, timeToPrepare, healthy, nil
}

//functions for getting and setting the
