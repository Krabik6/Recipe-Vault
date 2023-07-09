package meals

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// HandleMessage handles message for create meal state
func (cms *CreateMealStateHandler) HandleMessage(ctx context.Context, userID int64, update tgbotapi.Update) error {
	state, err := cms.GetUserState(ctx, userID)
	if err != nil {
		return err
	}
	cms.LocalState = state

	fmt.Println(cms.LocalState, "Create meal state")
	message := update.Message.Text
	if message == "/cancel" {
		return cms.handleCancel(ctx, userID)
	}
	switch cms.LocalState {
	//case NoCreateMealState:
	//	return cms.handleNoState(userID, update)
	case CreateMealName:
		cms.Name = message
	case CreateMealTime:
		cms.Time = message
	case CreateMealRecipes:
		if message == "confirm" {
			err := cms.handleCreateMealRecipes(userID)
			if err != nil {
				return err
			}
		} else {
			err := cms.handleState(ctx, userID, update)
			log.Println("here")
			if err != nil {
				return err
			}

		}

	case CreateMealConfirmation:
		if message == "yes" {
			return cms.handleCreateMealConfirmYes(ctx, userID)
		} else if message == "no" {
			return cms.handleCreateMealConfirmNo(ctx, userID)
		} else {
			// msg to user: press one of the buttons or /cancel or print yes or no to confirm
			reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
			msg := tgbotapi.NewMessage(update.Message.From.ID, reply)
			_, err := cms.Bot.Send(msg)
			if err != nil {
				return err
			}
			return fmt.Errorf("unknown state there: %d", cms.LocalState)
		}
	}

	err = cms.handleState(ctx, userID, update)
	if err != nil {
		return err
	}

	return nil
}
