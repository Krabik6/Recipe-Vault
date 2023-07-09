package recipes

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

// HandleMessage handles the message
func (crs *CreateRecipeService) HandleMessage(
	ctx context.Context,
	userID int64,
	update tgbotapi.Update,
) error {

	recipeCreationState, err := crs.GetUserState(ctx, userID)
	if err != nil {
		return err
	}
	crs.LocalState = recipeCreationState

	fmt.Printf("current state: %d\n", crs.LocalState)
	message := update.Message.Text
	if message == "/cancel" {
		return crs.handleCancel(ctx, userID)
	}
	switch crs.LocalState {
	case CreateRecipeTitle:
		crs.Title = message
	case CreateRecipeDescription:
		crs.Description = message
	case CreateRecipeIsPublic:
		if message == "yes" {
			crs.IsPublic = true
		} else {
			crs.IsPublic = false
		}
	case CreateRecipeCost:
		cost, err := strconv.ParseFloat(message, 64)
		if err != nil {
			return fmt.Errorf("error parsing cost: %v", err)
		}
		crs.Cost = cost
	case CreateRecipeTimeToPrepare:
		timeToPrepare, err := strconv.ParseInt(message, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing timeToPrepare: %v", err)
		}
		crs.TimeToPrepare = timeToPrepare
	case CreateRecipeHealthy:
		healthy, err := strconv.Atoi(message)
		if err != nil {
			return fmt.Errorf("error parsing healthy: %v", err)
		}
		crs.Healthy = healthy
	case CreateRecipeConfirmation:
		if message == "yes" {
			return crs.handleCreateRecipeConfirmYes(ctx, userID)
		} else if message == "no" {
			return crs.handleCreateRecipeConfirmNo(ctx, userID)
		} else {
			// msg to user: press one of the buttons or /cancel or print yes or no to confirm
			reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
			msg := tgbotapi.NewMessage(update.Message.From.ID, reply)
			_, err := crs.Bot.Send(msg)
			if err != nil {
				return err
			}
			return fmt.Errorf("unknown state: %d", crs.LocalState)
		}
	}

	//handle state
	err = crs.handleState(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
