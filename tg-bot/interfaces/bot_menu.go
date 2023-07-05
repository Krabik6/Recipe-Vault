package interfaces

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotMenu interface {
	CreateMainMenu(ctx context.Context, userID int64) tgbotapi.ReplyKeyboardMarkup
	CreateCancelKeyboard() tgbotapi.ReplyKeyboardMarkup
}
