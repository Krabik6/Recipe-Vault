package meals

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCancel that cancel the creation of a recipe: delete the state from manager and send a message to the user
func (cms *CreateMealStateHandler) handleCancel(ctx context.Context, userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user

	cms.LocalState = NoCreateMealState
	err := cms.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = cms.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = cms.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Creation of meal plan canceled")
	msg.ReplyMarkup = cms.BotMenu.CreateMainMenu(ctx, userID)
	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// handleNoCreateMealState handles no create meal state
func (cms *CreateMealStateHandler) handleNoCreateMealState(userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	msg := tgbotapi.NewMessage(userID, "Enter the name of the meal plan")
	_, err := cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealName
	return nil
}

// handleCreateMealName handles create meal name state
func (cms *CreateMealStateHandler) handleCreateMealName(userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	//format 2022-05-21 00:0:31
	msg := tgbotapi.NewMessage(userID, "Enter the time of the meal plan in the format 2022-05-21 00:0:31")
	_, err := cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealTime
	return nil
}

// handleCreateMealTime handles create meal time state
func (cms *CreateMealStateHandler) handleCreateMealTime(ctx context.Context, userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	//format 2022-05-21 00:0:31
	err := cms.RecipesService.RecipesList(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Press the button with the recipe number to add it to the meal plan")

	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("confirm"),
		tgbotapi.NewKeyboardButton("/cancel"),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1)
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard

	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealRecipes
	return nil
}

// handleCreateMealRecipes handles create meal recipes state, send a message to the user with the the info about the meal plan
func (cms *CreateMealStateHandler) handleCreateMealRecipes(userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	reply := fmt.Sprintf("Meal plan: \n Name: %s \n Time: %s \n Recipes: %v", cms.Name, cms.Time, cms.Recipes)
	msg := tgbotapi.NewMessage(userID, reply)
	yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yesButton, noButton))
	_, err := cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealConfirmation
	return nil
}

// handleCreateMealConfirmYes handles the CreateMealConfirmation state when the user confirms the creation of the meal: set the state to NoCreateMealState and send a message to the user that the meal is created
func (cms *CreateMealStateHandler) handleCreateMealConfirmYes(ctx context.Context, userID int64) error {
	meal := &models.Meal{
		Name:    cms.Name,
		AtTime:  cms.Time,
		Recipes: cms.Recipes,
	}
	token, err := cms.JwtManager.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	err = cms.Api.CreateMealPlan(*meal, token)
	if err != nil {
		return err
	}
	cms.LocalState = NoCreateMealState
	err = cms.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(userID, "Meal plan created on date: "+cms.Time)
	msg.ReplyMarkup = cms.BotMenu.CreateMainMenu(ctx, userID)
	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// handleCreateMealConfirmNo handles the CreateMealConfirmation state when the user does not confirm the creation of the meal: set the state to NoCreateMealState and send a message to the user that the meal is not created
func (cms *CreateMealStateHandler) handleCreateMealConfirmNo(ctx context.Context, userID int64) error {
	cms.LocalState = NoCreateMealState
	err := cms.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(userID, "Meal plan not created")
	msg.ReplyMarkup = cms.BotMenu.CreateMainMenu(ctx, userID)
	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
