package login

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// HandleCallbackQuery callback query states
func (ls *LoginStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	state, err := ls.GetUserLoginState(ctx, userID)
	if err != nil {
		return err
	}
	ls.LocalState = state
	query := update.CallbackQuery
	data := query.Data
	log.Println("state: ", ls.LocalState)
	switch ls.LocalState {
	case LoginConfirmation:
		if data == model.ConfirmButton {
			log.Println("Login confirmation")
			return ls.handleLoginComplete(ctx, userID, update)
		} else if data == model.CancelButton {
			message := tgbotapi.NewMessage(userID, "Вход в аккаунт отменен")
			_, err := ls.Bot.Send(message)
			if err != nil {
				return err
			}
			ls.LocalState = NoLoginState
			err = ls.SetUserLoginState(ctx, userID)
			if err != nil {
				return err
			}
			err = ls.DeleteUserLoginData(ctx, userID)
			if err != nil {
				return err
			}
			err = ls.UserStateManager.DeleteUserState(ctx, userID)
			if err != nil {
				return err
			}
			return nil
		} else {

			return fmt.Errorf("unknown command")
		}
	case NoLoginState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = ls.BotMenu.CreateCancelKeyboard()
		_, err := ls.Bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	default:
		// print query
		log.Println(query.Data)
		return fmt.Errorf("unknown login state")
	}
}
