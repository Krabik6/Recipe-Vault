package registration

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (rs *RegistrationStateHandler) handleNoRegistrationState(userID int64) error {
	fmt.Println("Welcome! Please provide your name.")
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Welcome! Please provide your name."))
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на имя
	rs.LocalState = RegistrationName
	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationName(userID int64) error {
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Hello, %s! Please provide your email.\n", rs.Name)))
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на email
	rs.LocalState = RegistrationEmail

	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationEmail(ctx context.Context, userID int64) error {
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Username %s is registered. Please provide your password.\n", rs.Email)))
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на пароль
	rs.LocalState = RegistrationPassword

	// Сохраняем состояние в Redis
	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationPassword(ctx context.Context, userID int64) error {
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Password is set. Please confirm your password."))
	if err != nil {
		return err
	}

	// Меняем состояние регистрации на подтверждение пароля
	rs.LocalState = RegistrationConfirmPassword

	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationConfirmPassword(ctx context.Context, userID int64, update tgbotapi.Update) error {
	if rs.ConfirmedPwd == rs.Password {
		_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Passwords match."))
		if err != nil {
			return err
		}
		rs.LocalState = RegistrationConfirmation
		err = rs.handleRegistrationConfirmation(ctx, userID, update)
		if err != nil {
			return err
		}
	} else {
		_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Password confirmation failed. Please try again. \n Provide your password."))
		if err != nil {
			return err
		}
		rs.LocalState = RegistrationPassword
	}

	// Сохраняем состояние в Redis
	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationConfirmation(ctx context.Context, userID int64, update tgbotapi.Update) error {
	reply := fmt.Sprintf("Please confirm your registration data.\nName: %s\nUsername: %s\nPassword: %s\n", rs.Name, rs.Email, rs.Password)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	// Кнопки подтверждения и отмены
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData(model.ConfirmButton, model.ConfirmButton)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData(model.CancelButton, model.CancelButton)
	// Добавляем кнопки в сообщение
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(confirmBtn, cancelBtn))
	// Отправляем сообщение пользователю
	_, err := rs.Bot.Send(msg)
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на ожидание подтверждения
	rs.LocalState = RegistrationConfirmation
	return nil

}

func (rs *RegistrationStateHandler) handleRegistrationComplete(ctx context.Context, userID int64, update tgbotapi.Update) error {

	user := model.SignUpCredentials{
		Username: rs.Email,
		Password: rs.Password,
		Name:     rs.Name,
	}
	err := rs.Api.SignUp(user)
	if err != nil {
		return err
	}
	// Обнуляем состояние регистрации
	rs.LocalState = NoRegistrationState

	//delete state from manager
	err = rs.DeleteUserRegistrationState(ctx, userID)
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

	msg := tgbotapi.NewMessage(userID, "Registration complete!\n")
	msg.ReplyMarkup = rs.BotMenu.CreateMainMenu(ctx, userID)
	_, err = rs.Bot.Send(msg)

	if err != nil {
		return err
	}

	// Удаляем состояние из Redis, так как регистрация завершена
	return nil
}
