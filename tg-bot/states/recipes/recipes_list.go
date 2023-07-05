package recipes

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"strconv"
)

// RecipesList функция для обработки команды /recipes_list в состоянии без состояния
func (sh *StateHandler) RecipesList(ctx context.Context, userID int64) error {
	client := &http.Client{}
	token, err := sh.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	recipes, err := api.GetRecipes(client, token)
	if err != nil {
		return err
	}

	//msg := tgbotapi.NewMessage(userID, "Recipes list:")
	for _, recipe := range recipes {
		msg := tgbotapi.NewMessage(userID, "")
		msg.Text += fmt.Sprintf("\n*Title*: %s", recipe.Title)
		msg.Text += fmt.Sprintf("\n*Description*: %s", recipe.Description)
		msg.Text += fmt.Sprintf("\n*Cost*: %.2f", recipe.Cost)
		msg.Text += fmt.Sprintf("\n*Time to prepare*: %d", recipe.TimeToPrepare)
		msg.Text += fmt.Sprintf("\n*Healthy(1-3)*: %d", recipe.Healthy)

		// Создаем CallbackData с ID рецепта
		callbackData := fmt.Sprintf(strconv.Itoa(recipe.Id))

		// Создаем инлайн-кнопку с текстом и CallbackData
		button := tgbotapi.NewInlineKeyboardButtonData(recipe.Title, callbackData)

		// Создаем клавиатуру с одной кнопкой и привязываем ее к сообщению
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(button),
		)
		msg.ReplyMarkup = inlineKeyboard
		msg.ParseMode = "Markdown"

		// Отправляем сообщение с кнопкой
		_, err := sh.Bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
