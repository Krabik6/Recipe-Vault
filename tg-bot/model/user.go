package model

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
)

type State int

const (
	NoState State = iota
	RegistrationState
	RecipeCreationState
	LogInState
	CreateMealState
	// Другие состояния
)

// constants for commands (start, registration, etc)
// constant for manager key (user state)
const UserState = "user_state:%d"

func (s State) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return nil, fmt.Errorf("failed to marshal RegistrationState: %v", err)
	}
	return data, nil
}

func (s *State) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, s)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return fmt.Errorf("failed to unmarshal RegistrationState: %v", err)
	}
	return nil
}
