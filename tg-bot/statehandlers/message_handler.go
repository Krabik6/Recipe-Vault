package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime/debug"
)

// HandleMessage Метод для обработки входящих сообщений в соответствии с текущим состоянием
func (sh *StateHandler) HandleMessage(
	ctx context.Context,
	userID int64,
	message string,
	update tgbotapi.Update,
) error {
	//get user state from manager
	state, err := sh.UserStateManager.GetUserState(ctx, userID)
	if err != nil {
		return err
	}

	switch state {
	case model.NoState:
		noStateHandler := sh.NoState
		state, err = noStateHandler.HandleMessage(ctx, userID, message, state)
		if err != nil {
			return err
		}
		if state != model.NoState {
			return sh.HandleMessage(ctx, userID, message, update)
		}
		return nil
	case model.RegistrationState:
		registrationHandler := sh.Registration
		name, email, password, err := registrationHandler.GetUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		registrationHandler.BuildUserData(name, email, password)
		return registrationHandler.HandleMessage(ctx, userID, message, update)
	case model.RecipeCreationState:
		recipeCreationHandler := sh.Recipes
		// set recipe creation data from manager
		title, description, isPublic, cost, timeToPrepare, healthy, err := recipeCreationHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.BuildUserData(title, description, isPublic, cost, timeToPrepare, healthy)

		return recipeCreationHandler.HandleMessage(ctx, userID, update)
	case model.CreateMealState:
		createMealHandler := sh.Meals
		// set recipe creation data from manager
		name, time, recipes, err := createMealHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}

		createMealHandler.BuildUserData(name, time, recipes)

		return createMealHandler.HandleMessage(ctx, userID, update)

	case model.LogInState:
		logInHandler := sh.Login
		username, password, err := logInHandler.GetUserLoginData(ctx, userID)
		if err != nil {
			return err
		}

		logInHandler.BuildUserData(username, password)
		return logInHandler.HandleMessage(ctx, userID, message, update)
	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}
}
