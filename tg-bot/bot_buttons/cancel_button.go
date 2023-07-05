package bot_buttons

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const cancelCommand = "/cancel"

func (bm *BotMenu) CreateCancelKeyboard() tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(cancelCommand),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1)
	keyboard.OneTimeKeyboard = true
	return keyboard
}
