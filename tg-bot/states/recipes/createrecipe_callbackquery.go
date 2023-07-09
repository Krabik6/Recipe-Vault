package recipes

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handle CallbackQuery
func (crs *CreateRecipeService) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	recipeCreationState, err := crs.GetUserState(ctx, userID)
	if err != nil {
		return err
	}
	crs.LocalState = recipeCreationState

	fmt.Printf("current state: %d\n", crs.LocalState)

	query := update.CallbackQuery
	data := query.Data
	switch crs.LocalState {
	case CreateRecipeIsPublic:
		if data == "yes" {
			crs.IsPublic = true
		} else {
			crs.IsPublic = false
		}
	case CreateRecipeHealthy:
		switch data {
		case "1":
			crs.Healthy = 1
		case "2":
			crs.Healthy = 2
		case "3":
			crs.Healthy = 3
		default:
			return fmt.Errorf("unknown healthy: %s", data)
		}
	case CreateRecipeConfirmation:
		if data == "yes" {
			return crs.handleCreateRecipeConfirmYes(ctx, userID)
		} else {
			return crs.handleCreateRecipeConfirmNo(ctx, userID)
		}
	default:
		// msg to user: press one of the buttons or /cancel or print yes or no to confirm
		reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
		msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, reply)
		_, err := crs.Bot.Send(msg)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown state: %d", crs.LocalState)
	}

	err = crs.handleState(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
