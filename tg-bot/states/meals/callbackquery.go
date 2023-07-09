package meals

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (cms *CreateMealStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	//get and set state
	state, err := cms.GetUserState(ctx, userID)
	if err != nil {
		return err
	}
	cms.LocalState = state

	fmt.Println(cms.LocalState, "Create meal state")
	query := update.CallbackQuery
	data := query.Data
	switch cms.LocalState {
	case CreateMealConfirmation:
		if data == "yes" {
			return cms.handleCreateMealConfirmYes(ctx, userID)
		} else {
			return cms.handleCreateMealConfirmNo(ctx, userID)
		}
	case CreateMealRecipes:
		if data == "confirm" {
			err := cms.handleCreateMealRecipes(userID)
			if err != nil {
				return err
			}
		} else if data == "yes" {
			return cms.handleCreateMealConfirmYes(ctx, userID)
		} else if data == "no" {
			return cms.handleCreateMealConfirmNo(ctx, userID)
		} else {
			err := cms.handleState(ctx, userID, update)
			log.Println("here")
			if err != nil {
				return err
			}

		}
	default:
		// msg to user: press one of the buttons or /cancel or print yes or no to confirm
		reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
		msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, reply)
		_, err := cms.Bot.Send(msg)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown state: damn%d", cms.LocalState)
	}

	//err := cms.handleState(ctx, userID)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
	return nil
}
