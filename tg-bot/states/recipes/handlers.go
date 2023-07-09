package recipes

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/context"
)

// handleNoCreateRecipeState handles the NoCreateRecipeState state: set the state to CreateRecipeTitle and send a message to the user to ask for the title
func (crs *CreateRecipeService) handleNoCreateRecipeState(userID int64) error {
	// set the state to CreateRecipeTitle and send a message to the user to ask for the title

	msg := tgbotapi.NewMessage(userID, "What is the title of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeTitle

	return nil
}

// handleCreateRecipeTitle handles the CreateRecipeTitle state: set the state to CreateRecipeDescription and send a message to the user to ask for the description
func (crs *CreateRecipeService) handleCreateRecipeTitle(userID int64) error {
	// set the state to CreateRecipeDescription and send a message to the user to ask for the description

	msg := tgbotapi.NewMessage(userID, "What is the description of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeDescription

	return nil
}

// handleCreateRecipeDescription handles the CreateRecipeDescription state: set the state to CreateRecipeIsPublic and send a message to the user to ask if the recipe is public
func (crs *CreateRecipeService) handleCreateRecipeDescription(userID int64) error {
	// set the state to CreateRecipeIsPublic and send a message with the keyboard to the user to ask if the recipe is public
	reply := fmt.Sprintf("Is the recipe public?")
	msg := tgbotapi.NewMessage(userID, reply)
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("no", "no")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(confirmBtn, cancelBtn))
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeIsPublic

	return nil
}

// handleCreateRecipeIsPublic handles the CreateRecipeIsPublic state: set the state to CreateRecipeCost and send a message to the user to ask for the cost
func (crs *CreateRecipeService) handleCreateRecipeIsPublic(userID int64) error {
	// set the state to CreateRecipeCost and send a message to the user to ask for the cost

	msg := tgbotapi.NewMessage(userID, "What is the cost of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeCost

	return nil
}

// handleCreateRecipeCost handles the CreateRecipeCost state: set the state to CreateRecipeTimeToPrepare and send a message to the user to ask for the time to prepare
func (crs *CreateRecipeService) handleCreateRecipeCost(userID int64) error {
	// set the state to CreateRecipeTimeToPrepare and send a message to the user to ask for the time to prepare

	msg := tgbotapi.NewMessage(userID, "What is the time to prepare of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeTimeToPrepare

	return nil
}

// handleCreateRecipeTimeToPrepare handles the CreateRecipeTimeToPrepare state: set the state to CreateRecipeHealthy and send a message to the user to ask if the recipe is healthy
func (crs *CreateRecipeService) handleCreateRecipeTimeToPrepare(userID int64) error {
	// set the state to CreateRecipeHealthy and send a message with the 3 buttons to the user to ask how healthy is the recipe (not healthy, healthy, very healthy)
	reply := fmt.Sprintf("How healthy is the recipe?")
	msg := tgbotapi.NewMessage(userID, reply)
	notHealthyBtn := tgbotapi.NewInlineKeyboardButtonData("not healthy", "1")
	healthyBtn := tgbotapi.NewInlineKeyboardButtonData("healthy", "2")
	veryHealthyBtn := tgbotapi.NewInlineKeyboardButtonData("very healthy", "3")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(notHealthyBtn, healthyBtn, veryHealthyBtn))

	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	crs.LocalState = CreateRecipeHealthy

	return nil
}

// handleCreateRecipeHealthy handles the CreateRecipeHealthy state: set the state to CreateRecipeConfirmation and send a message to the user to ask for the confirmation
func (crs *CreateRecipeService) handleCreateRecipeHealthy(userID int64) error {
	// set the state to CreateRecipeConfirmation and send a message with info about recipe and with the 2 buttons to the user to ask for the confirmation (yes, no)
	reply := fmt.Sprintf("Recipe info:\nTitle: %s\nDescription: %s\nIs public: %t\nCost: %d\nTime to prepare: %d\nHealthy: %s\n\nIs the info correct?", crs.Title, crs.Description, crs.IsPublic, crs.Cost, crs.TimeToPrepare, crs.Healthy)
	// create repy with beutiful format of the recipe info
	reply = fmt.Sprintf("Recipe info:\nTitle: %s\nDescription: %s\nIs public: %t\nCost: %.2f\nTime to prepare: %d\nHealthy: %d\n\nIs the info correct?", crs.Title, crs.Description, crs.IsPublic, crs.Cost, crs.TimeToPrepare, crs.Healthy)
	msg := tgbotapi.NewMessage(userID, reply)
	yesBtn := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noBtn := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yesBtn, noBtn))

	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	crs.LocalState = CreateRecipeConfirmation

	return nil
}

// handleCreateRecipeConfirmYes handles the CreateRecipeConfirmation state when the user confirms the creation of the recipe: set the state to NoCreateRecipeState and send a message to the user that the recipe is created
func (crs *CreateRecipeService) handleCreateRecipeConfirmYes(ctx context.Context, userID int64) error {
	// set the state to NoCreateRecipeState and send a message to the user that the recipe is created

	// create the recipe
	recipe := model.CreateRecipeInput{
		Title:         crs.Title,
		Description:   crs.Description,
		IsPublic:      crs.IsPublic,
		Cost:          crs.Cost,
		TimeToPrepare: crs.TimeToPrepare,
		Healthy:       crs.Healthy,
	}
	//gwt jwt token
	token, err := crs.JwtManager.GetUserJWTToken(ctx, userID)

	recipeId, err := crs.Api.CreateRecipe(recipe, token)
	if err != nil {
		return err
	}

	crs.LocalState = NoCreateRecipeState
	err = crs.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Recipe created with id %d", recipeId))
	msg.ReplyMarkup = crs.BotMenu.CreateMainMenu(ctx, userID)
	_, err = crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// handleCreateRecipeConfirmNo handles the CreateRecipeConfirmation state when the user doesn't confirm the creation of the recipe: set the state to NoCreateRecipeState and send a message to the user that the recipe is not created also delete
func (crs *CreateRecipeService) handleCreateRecipeConfirmNo(ctx context.Context, userID int64) error {
	// set the state to NoCreateRecipeState and send a message to the user that the recipe is not created

	crs.LocalState = NoCreateRecipeState
	err := crs.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Creation of recipe canceled")
	_, err = crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// handleCancel that cancel the creation of a recipe: delete the state from manager and send a message to the user
func (crs *CreateRecipeService) handleCancel(ctx context.Context, userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user

	crs.LocalState = NoCreateRecipeState
	err := crs.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Creation of recipe canceled")
	msg.ReplyMarkup = crs.BotMenu.CreateMainMenu(ctx, userID)
	_, err = crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
