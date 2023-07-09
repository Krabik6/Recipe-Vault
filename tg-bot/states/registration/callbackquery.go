package registration

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// HandleCallbackQuery callback query states
func (rs *RegistrationStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	//get and set user state
	state, err := rs.GetUserRegistrationState(ctx, userID)
	if err != nil {
		return err
	}
	rs.LocalState = state

	query := update.CallbackQuery
	data := query.Data
	log.Println("state: ", rs.LocalState)
	switch rs.LocalState {
	case RegistrationConfirmation:
		if data == model.ConfirmButton {
			log.Println("подтверждение регистрации")

			return rs.handleRegistrationComplete(ctx, userID, update)

		} else if data == model.CancelButton {
			message := tgbotapi.NewMessage(userID, "Регистрация отменена")
			_, err := rs.Bot.Send(message)
			if err != nil {
				return err
			}

			// Отмена регистрации, переход к начальному состоянию
			rs.LocalState = NoRegistrationState
			err = rs.SetUserRegistrationState(ctx, userID)
			if err != nil {
				return err
			}
			err = rs.DeleteUserRegistrationData(ctx, userID)
			if err != nil {
				return err
			}
			err = rs.UserStateManager.DeleteUserState(ctx, userID)
			if err != nil {
				return err
			}

			return nil
		} else {
			return fmt.Errorf("unknown command")
		}
	default:
		// print query
		log.Println(query.Data)
		return fmt.Errorf("unknown registration state")
	}
}
