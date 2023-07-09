package recipes

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type ListService struct {
	JwtManager interfaces.JwtManager
	Bot        *tgbotapi.BotAPI
	Api        *api.Api
}

func NewListService(
	jwtManager interfaces.JwtManager,
	bot *tgbotapi.BotAPI,
	api *api.Api,
) *ListService {
	return &ListService{
		JwtManager: jwtManager,
		Bot:        bot,
		Api:        api,
	}
}

// RecipesList функция для обработки команды /recipes_list в состоянии без состояния
func (ls *ListService) RecipesList(ctx context.Context, userID int64) error {
	token, err := ls.JwtManager.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	recipes, err := ls.Api.GetRecipes(token)
	if err != nil {
		return err
	}

	//msg := tgbotapi.NewMessage(userID, "RecipesService list:")
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
		_, err := ls.Bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
