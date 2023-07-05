package statehandlers

// file name: loginstatehandler_utils.go
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
)

func (ls loginState) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(ls)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return nil, fmt.Errorf("failed to marshal LoginState: %v", err)
	}
	return data, nil
}

func (ls *loginState) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, ls)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return fmt.Errorf("failed to unmarshal LoginState: %v", err)
	}
	return nil
}

// SetUserLoginState устанавливает состояние входа пользователя в Redis.
func (ls *LoginStateHandler) SetUserLoginState(ctx context.Context, userID int64) error {
	state, err := ls.LocalState.MarshalBinary()
	if err != nil {
		return err
	}

	key := fmt.Sprintf(userLoginState, userID)

	err = ls.Client.Set(ctx, key, state, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserLoginState удаляет состояние входа пользователя из Redis.
func (ls *LoginStateHandler) DeleteUserLoginState(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(userLoginState, userID)
	err := ls.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserLoginState получает состояние входа пользователя из Redis.
func (ls *LoginStateHandler) GetUserLoginState(ctx context.Context, userID int64) (loginState, error) {
	key := fmt.Sprintf(userLoginState, userID)
	stateStr, err := ls.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return NoLoginState, nil
		}
		return NoLoginState, err
	}

	var state loginState
	err = state.UnmarshalBinary([]byte(stateStr))
	if err != nil {
		return NoLoginState, err
	}

	return state, nil
}

// SetUserLoginData устанавливает данные входа пользователя в Redis.
func (ls *LoginStateHandler) SetUserLoginData(ctx context.Context, userID int64) error {
	keyEmail := fmt.Sprintf(userLoginEmail, userID)
	err := ls.Client.Set(ctx, keyEmail, ls.Username, 0).Err()
	if err != nil {
		return err
	}

	keyPassword := fmt.Sprintf(userLoginPassword, userID)
	err = ls.Client.Set(ctx, keyPassword, ls.Password, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserLoginData удаляет данные входа пользователя из Redis.
func (ls *LoginStateHandler) DeleteUserLoginData(ctx context.Context, userID int64) error {
	keyEmail := fmt.Sprintf(userLoginEmail, userID)
	err := ls.Client.Del(ctx, keyEmail).Err()
	if err != nil {
		return err
	}

	keyPassword := fmt.Sprintf(userLoginPassword, userID)
	err = ls.Client.Del(ctx, keyPassword).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserLoginData получает данные входа пользователя из Redis.
func (ls *LoginStateHandler) GetUserLoginData(ctx context.Context, userID int64) (string, string, error) {
	keyEmail := fmt.Sprintf(userLoginEmail, userID)
	email, err := ls.Client.Get(ctx, keyEmail).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", nil
		}
		return "", "", err
	}

	keyPassword := fmt.Sprintf(userLoginPassword, userID)
	password, err := ls.Client.Get(ctx, keyPassword).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", nil
		}
		return "", "", err
	}

	return email, password, nil
}
