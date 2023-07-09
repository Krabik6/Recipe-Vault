package bot_buttons

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/manager"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BotMenu struct {
	Bot        *tgbotapi.BotAPI
	JwtManager *manager.JwtManager
}

func NewBotMenu(bot *tgbotapi.BotAPI, jwtManager *manager.JwtManager) *BotMenu {
	return &BotMenu{Bot: bot, JwtManager: jwtManager}
}
func (bm *BotMenu) CreateMainMenu(ctx context.Context, userID int64) tgbotapi.ReplyKeyboardMarkup {
	loggedIn, err := bm.JwtManager.CheckLoggedIn(ctx, userID)
	log.Println("current login status:", loggedIn)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}
	}

	if loggedIn {
		return jwtMenu()
	} else {
		return noJwtMenu()
	}
}

func jwtMenu() tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(model.HelpCommand),
		tgbotapi.NewKeyboardButton(model.LogOutCommand),
	}

	row2 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(model.CreateRecipeCommand),
		tgbotapi.NewKeyboardButton(model.RecipesListCommand),
		tgbotapi.NewKeyboardButton(model.CreateMealCommand),
		tgbotapi.NewKeyboardButton(model.MealsListCommand),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	keyboard.OneTimeKeyboard = false // Здесь изменено значение на false
	return keyboard
}

func noJwtMenu() tgbotapi.ReplyKeyboardMarkup {
	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(model.HelpCommand),
	}

	row2 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(model.RegistrationCommand),
		tgbotapi.NewKeyboardButton(model.LogInCommand),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2)
	keyboard.OneTimeKeyboard = false // Здесь изменено значение на false
	return keyboard
}
