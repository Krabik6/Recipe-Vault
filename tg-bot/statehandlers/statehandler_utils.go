package statehandlers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func (s State) MarshalBinary() ([]byte, error) {
	// Преобразование значения State в байтовый массив
	data := []byte(strconv.Itoa(int(s)))
	return data, nil
}

func (s *State) UnmarshalBinary(data []byte) error {
	// Преобразование байтового массива в значение State
	value, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	*s = State(value)
	return nil
}

// setUserState sets the user state in Redis.
func (sh *StateHandler) setUserState(ctx context.Context, userID int64) error {
	// Forming the key for the user
	key := fmt.Sprintf(userState, userID)

	// Setting the state value in Redis
	err := sh.Client.Set(ctx, key, sh.State, 0).Err()
	if err != nil {
		return err // Return the error if there was a problem setting the state in Redis
	}

	return nil // Return nil to indicate success
}

// getUserState retrieves the user state from Redis.
// If the state is not found, it returns NoState.
func (sh *StateHandler) getUserState(ctx context.Context, userID int64) (State, error) {
	// Forming the key for the user
	key := fmt.Sprintf(userState, userID)

	// Getting the state value from Redis
	stateStr, err := sh.Client.Get(ctx, key).Result()
	if err != nil {
		// Handling the case when the value is not found in Redis
		if err == redis.Nil {
			return NoState, nil // Return NoState if the state is not found
		}
		return NoState, err // Return an error if there is any other Redis error
	}

	var state State
	err = state.UnmarshalBinary([]byte(stateStr)) // Deserializing the state
	if err != nil {
		return NoState, err // Return an error if there is an error in deserializing the state
	}

	sh.State = state // Set the state in the StateHandler struct

	return state, nil // Return the user state
}

// deleteUserState deletes the user state from Redis.
func (sh *StateHandler) deleteUserState(ctx context.Context, userID int64) error {
	// Forming the key for the user
	key := fmt.Sprintf(userState, userID)

	// Deleting the state value from Redis
	err := sh.Client.Del(ctx, key).Err()
	if err != nil {
		return err // Return the error if there was a problem deleting the state from Redis
	}

	return nil // Return nil to indicate success
}
