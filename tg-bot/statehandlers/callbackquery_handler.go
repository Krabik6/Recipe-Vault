package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"runtime/debug"
)

// HandleCallbackQuery handle callback query
func (sh *StateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	state, err := sh.UserStateManager.GetUserState(ctx, userID)
	if err != nil {
		return err
	}
	switch state {
	case model.RecipeCreationState:
		recipeCreationHandler := sh.Recipes
		// set recipe creation data from manager
		title, description, isPublic, cost, timeToPrepare, healthy, err := recipeCreationHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.BuildUserData(title, description, isPublic, cost, timeToPrepare, healthy)

		return recipeCreationHandler.HandleCallbackQuery(ctx, userID, update)
	case model.RegistrationState:
		rsh := sh.Registration
		name, email, password, err := rsh.GetUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		rsh.BuildUserData(name, email, password)
		err = rsh.HandleCallbackQuery(ctx, userID, update)
		if err != nil {
			return err
		}
	case model.LogInState:
		logInHandler := sh.Login
		username, password, err := logInHandler.GetUserLoginData(ctx, userID)
		if err != nil {
			return err
		}

		logInHandler.BuildUserData(username, password)

		return logInHandler.HandleCallbackQuery(ctx, userID, update)
	case model.CreateMealState:
		createMealHandler := sh.Meals
		// set recipe creation data from manager
		name, time, recipes, err := createMealHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}

		createMealHandler.BuildUserData(name, time, recipes)
		return createMealHandler.HandleCallbackQuery(ctx, userID, update)

	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}

	return nil
}
