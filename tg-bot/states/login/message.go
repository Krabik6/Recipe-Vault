package login

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
)

func (ls *LoginStateHandler) HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error {
	state, err := ls.GetUserLoginState(ctx, userID)
	if err != nil {
		return err
	}
	ls.LocalState = state
	// Check if the message is "/cancel" to cancel the login process
	if message == "/cancel" {
		ls.LocalState = NoLoginState
		err := ls.SetUserLoginState(ctx, userID)
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
		// Add message to the user
		msg := tgbotapi.NewMessage(userID, "Login canceled")
		msg.ReplyMarkup = ls.BotMenu.CreateMainMenu(ctx, userID)
		_, err = ls.Bot.Send(msg)
		if err != nil {
			return err
		}
		log.Println("Login canceled, state: ", ls.LocalState)
		return nil
	}
	switch ls.LocalState {
	case NoLoginState:
		//send cancel button without message
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = ls.BotMenu.CreateCancelKeyboard()
		_, err := ls.Bot.Send(msg)
		if err != nil {
			return err
		}
	case LoginEmail:
		ls.Username = message
	case LoginPassword:
		ls.Password = message
	}

	// Handle state
	err = ls.HandleState(ctx, userID, update)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return err
	}

	err = ls.SetUserLoginState(ctx, userID)
	if err != nil {
		return err
	}
	err = ls.SetUserLoginData(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
