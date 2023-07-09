package meals

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/states/recipes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
)

type createMealState int

// Name    string
// AtTime  string
// Recipes []int
const (
	NoCreateMealState createMealState = iota
	CreateMealName
	CreateMealTime
	CreateMealRecipes
	CreateMealConfirmation
)

type CreateMealStateHandler struct {
	Bot              *tgbotapi.BotAPI
	UserStateManager interfaces.UserStateManager
	JwtManager       interfaces.JwtManager
	Client           *redis.Client
	BotMenu          interfaces.BotMenu
	RecipesService   *recipes.RecipesService
	Api              *api.Api
	LocalState       createMealState
	Name             string
	Time             string
	Recipes          []int
}

func NewCreateMealStateHandler(
	bot *tgbotapi.BotAPI,
	userStateManager interfaces.UserStateManager,
	jwtManager interfaces.JwtManager,
	client *redis.Client,
	botMenu interfaces.BotMenu,
	recipesService *recipes.RecipesService,
	api *api.Api,
) *CreateMealStateHandler {
	return &CreateMealStateHandler{
		Bot:              bot,
		UserStateManager: userStateManager,
		JwtManager:       jwtManager,
		Client:           client,
		BotMenu:          botMenu,
		RecipesService:   recipesService,
		Api:              api,
	}
}

func (cms *CreateMealStateHandler) BuildUserData(name string, time string, recipes []int) {
	cms.Name = name
	cms.Time = time
	cms.Recipes = recipes
}

const (
	createMealStateKey   = "create_meal_state:%d"
	createMealNameKey    = "create_meal_name:%d"
	createMealTimeKey    = "create_meal_time:%d"
	createMealRecipesKey = "create_meal_recipes:%d"
)

func extractRecipeID(data string) (int, error) {
	// Удаляем префикс "view_recipe:"
	recipeIDStr := strings.TrimPrefix(data, "view_recipe:")

	// Преобразуем полученную строку в число
	recipeID, err := strconv.Atoi(recipeIDStr)
	if err != nil {
		return 0, err
	}

	return recipeID, nil
}

// handleState handles state for create meal state
func (cms *CreateMealStateHandler) handleState(ctx context.Context, userID int64, update tgbotapi.Update) error {
	fmt.Println(cms.LocalState, "Create meal state")
	switch cms.LocalState {
	case NoCreateMealState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = cms.BotMenu.CreateCancelKeyboard()
		_, err := cms.Bot.Send(msg)
		if err != nil {
			return err
		}
		err = cms.handleNoCreateMealState(userID)
		if err != nil {
			return err
		}
	case CreateMealName:
		err := cms.handleCreateMealName(userID)
		if err != nil {
			return err
		}
	case CreateMealTime:
		err := cms.handleCreateMealTime(ctx, userID)
		if err != nil {
			return err
		}

	case CreateMealRecipes:
		if update.CallbackQuery != nil {

			recipeID, err := extractRecipeID(update.CallbackQuery.Data)
			if err != nil {
				// Обработка ошибки
			}
			cms.Recipes = append(cms.Recipes, recipeID)
		}

		//print what recipes are addded already
		msg := tgbotapi.NewMessage(userID, fmt.Sprintf("You already added these recipes: %v", cms.Recipes))
		_, err := cms.Bot.Send(msg)
		if err != nil {
			return err
		}
	//case CreateMealConfirmation:
	//	err := cms.handleCreateMealConfirmation(userID)
	//	if err != nil {
	//		return err
	//	}
	default:
		return fmt.Errorf("unknown state: %d", cms.LocalState)
	}

	err := cms.SetUserData(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.SetUserState(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
