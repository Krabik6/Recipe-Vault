package expectation

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// LogOut функция для обработки команды /logout без состояния
func (nsh *NoStateHandler) LogOut(ctx context.Context, userID int64) error {
	//if user has already logged in - log out, else - send message that user is not logged in
	// Получение состояния пользователя
	token, err := nsh.JwtManager.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	if token != "" {
		// Удаление токена из базы данных
		err = nsh.JwtManager.DeleteUserJWTToken(ctx, userID)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(userID, "You logout successfully.")

		msg.ReplyMarkup = nsh.BotMenu.CreateMainMenu(ctx, userID)
		// Вывод сообщения о том, что пользователь вышел из аккаунта
		_, err := nsh.Bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	} else {
		// Вывод сообщения о том, что пользователь не зарегистрирован
		_, err := nsh.Bot.Send(tgbotapi.NewMessage(userID, "You are not logged in."))
		if err != nil {
			return err
		}
		return nil
	}
}
