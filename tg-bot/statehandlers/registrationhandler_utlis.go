package statehandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
)

func (rs registrationState) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(rs)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return nil, fmt.Errorf("failed to marshal RegistrationState: %v", err)
	}
	return data, nil
}

func (rs *registrationState) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, rs)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return fmt.Errorf("failed to unmarshal RegistrationState: %v", err)
	}
	return nil
}

// SetUserRegistrationState устанавливает состояние регистрации пользователя в Redis.
func (rs *RegistrationStateHandler) SetUserRegistrationState(ctx context.Context, userID int64) error {
	state, err := rs.State.MarshalBinary()
	if err != nil {
		return err
	}

	key := fmt.Sprintf(userRegistrationState, userID)

	err = rs.Client.Set(ctx, key, state, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserRegistrationState удаляет состояние регистрации пользователя из Redis.
func (rs *RegistrationStateHandler) DeleteUserRegistrationState(ctx context.Context, userID int64) error {
	key := fmt.Sprintf(userRegistrationState, userID)
	err := rs.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserRegistrationState получает состояние регистрации пользователя из Redis.
func (rs *RegistrationStateHandler) GetUserRegistrationState(ctx context.Context, userID int64) (registrationState, error) {
	key := fmt.Sprintf(userRegistrationState, userID)
	stateStr, err := rs.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return NoRegistrationState, nil
		}
		return NoRegistrationState, err
	}

	var state registrationState
	err = state.UnmarshalBinary([]byte(stateStr))
	if err != nil {
		return NoRegistrationState, err
	}

	return state, nil
}

// SetUserRegistrationData устанавливает данные регистрации пользователя в Redis.
func (rs *RegistrationStateHandler) SetUserRegistrationData(ctx context.Context, userID int64) error {
	keyName := fmt.Sprintf(userRegistrationName, userID)
	err := rs.Client.Set(ctx, keyName, rs.Name, 0).Err()
	if err != nil {
		return err
	}

	keyEmail := fmt.Sprintf(userRegistrationEmail, userID)
	err = rs.Client.Set(ctx, keyEmail, rs.Email, 0).Err()
	if err != nil {
		return err
	}

	keyPassword := fmt.Sprintf(userRegistrationPassword, userID)
	err = rs.Client.Set(ctx, keyPassword, rs.Password, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserRegistrationData удаляет данные регистрации пользователя из Redis.
func (rs *RegistrationStateHandler) DeleteUserRegistrationData(ctx context.Context, userID int64) error {
	keyName := fmt.Sprintf(userRegistrationName, userID)
	err := rs.Client.Del(ctx, keyName).Err()
	if err != nil {
		return err
	}

	keyEmail := fmt.Sprintf(userRegistrationEmail, userID)
	err = rs.Client.Del(ctx, keyEmail).Err()
	if err != nil {
		return err
	}

	keyPassword := fmt.Sprintf(userRegistrationPassword, userID)
	err = rs.Client.Del(ctx, keyPassword).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserRegistrationData получает данные регистрации пользователя из Redis.
func (rs *RegistrationStateHandler) GetUserRegistrationData(ctx context.Context, userID int64) (string, string, string, error) {
	keyName := fmt.Sprintf(userRegistrationName, userID)
	name, err := rs.Client.Get(ctx, keyName).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	keyEmail := fmt.Sprintf(userRegistrationEmail, userID)
	email, err := rs.Client.Get(ctx, keyEmail).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	keyPassword := fmt.Sprintf(userRegistrationPassword, userID)
	password, err := rs.Client.Get(ctx, keyPassword).Result()
	if err != nil {
		if err == redis.Nil {
			return "", "", "", nil
		}
		return "", "", "", err
	}

	return name, email, password, nil
}
