package login

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (ls *LoginStateHandler) handleNoLoginState(userID int64) error {
	// Send message to the user
	_, err := ls.Bot.Send(tgbotapi.NewMessage(userID, "Welcome! Please enter your username."))
	if err != nil {
		return err
	}
	// Change login state to username
	ls.LocalState = LoginEmail
	return nil
}

func (ls *LoginStateHandler) handleLoginEmail(userID int64) error {
	// Send message to the user
	_, err := ls.Bot.Send(tgbotapi.NewMessage(userID, "Please enter your email."))
	if err != nil {
		return err
	}
	// Change login state to password
	ls.LocalState = LoginEmail
	return nil
}

func (ls *LoginStateHandler) handleLoginPassword(userID int64) error {
	// Send message to the user
	_, err := ls.Bot.Send(tgbotapi.NewMessage(userID, "Email accepted. Please enter your password. "))
	if err != nil {
		return err
	}
	// Change login state to confirmation
	ls.LocalState = LoginPassword
	return nil
}

// handleLoginConfirmation handles the login confirmation
func (ls *LoginStateHandler) handleLoginConfirmation(ctx context.Context, userID int64, update tgbotapi.Update) error {
	// Send message to the user like above, but with login data instead of registration data
	reply := fmt.Sprintf("Please confirm your login data.\nUsername: %s\nPassword: %s\n", ls.Username, ls.Password)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	// Confirm and cancel buttons
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData(model.ConfirmButton, model.ConfirmButton)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData(model.CancelButton, model.CancelButton)
	// Add buttons to the message
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(confirmBtn, cancelBtn))
	// Send message to the user
	_, err := ls.Bot.Send(msg)
	if err != nil {
		return err
	}
	// Change login state to confirmation
	ls.LocalState = LoginConfirmation
	return nil
}

func (ls *LoginStateHandler) handleLoginComplete(ctx context.Context, userID int64, update tgbotapi.Update) error {
	user := model.LoginCredentials{
		Username: ls.Username,
		Password: ls.Password,
	}
	token, err := ls.Api.Login(user)
	if err != nil {
		return fmt.Errorf("error logging in: %v", err)
	}
	err = ls.JwtManager.SetUserJWTToken(ctx, userID, token)
	if err != nil {
		return err
	}
	// Reset login state
	ls.LocalState = NoLoginState

	// Delete state from Redis
	err = ls.DeleteUserLoginState(ctx, userID)
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

	// Add message to the user with the token
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Login complete! Your token is: %s\n", token))
	msg.ReplyMarkup = ls.BotMenu.CreateMainMenu(ctx, userID)

	_, err = ls.Bot.Send(msg)
	if err != nil {
		return err
	}
	// Remove state from Redis as login is completed
	return nil
}
