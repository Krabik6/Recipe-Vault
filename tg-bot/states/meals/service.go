package meals

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/states/recipes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type CreateMeal interface {
	SetUserState(ctx context.Context, userID int64) error
	DeleteUserState(ctx context.Context, userID int64) error
	GetUserState(ctx context.Context, userID int64) (createMealState, error)
	SetUserData(ctx context.Context, userID int64) error
	DeleteUserData(ctx context.Context, userID int64) error
	GetUserData(ctx context.Context, userID int64) (name string, time string, recipes []int, err error)
	BuildUserData(name string, time string, recipes []int)
	HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error
	HandleMessage(ctx context.Context, userID int64, update tgbotapi.Update) error
}

type List interface {
	MealPlansList(ctx context.Context, userID int64) error
}

type MealsService struct {
	CreateMeal
	List
}

func NewMealsService(
	bot *tgbotapi.BotAPI,
	userStateManager interfaces.UserStateManager,
	jwtManager interfaces.JwtManager,
	client *redis.Client,
	botMenu interfaces.BotMenu,
	recipesService *recipes.RecipesService,
	api *api.Api,
) *MealsService {
	return &MealsService{
		CreateMeal: NewCreateMealStateHandler(
			bot,
			userStateManager,
			jwtManager,
			client,
			botMenu,
			recipesService,
			api,
		),
		List: NewListService(jwtManager, bot, api),
	}
}
