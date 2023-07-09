package recipes

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type CreateRecipe interface {
	SetUserState(ctx context.Context, userID int64) error
	DeleteUserState(ctx context.Context, userID int64) error
	GetUserState(ctx context.Context, userID int64) (createRecipeState, error)
	SetUserData(ctx context.Context, userID int64) error
	DeleteUserData(ctx context.Context, userID int64) error
	GetUserData(ctx context.Context, userID int64) (title string, description string, isPublic bool, cost float64, timeToPrepare int64, healthy int, err error)
	HandleMessage(ctx context.Context, userID int64, update tgbotapi.Update) error
	HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error
	BuildUserData(title string, description string, isPublic bool, cost float64, timeToPrepare int64, healthy int)
}

type List interface {
	RecipesList(ctx context.Context, userID int64) error
}

type RecipesService struct {
	CreateRecipe
	List
}

func NewRecipesService(
	bot *tgbotapi.BotAPI,
	client *redis.Client,
	userManager interfaces.UserStateManager,
	jwtManager interfaces.JwtManager,
	botMenu interfaces.BotMenu,
	api *api.Api,
) *RecipesService {

	return &RecipesService{
		CreateRecipe: NewCreateRecipeStateHandler(bot, client, userManager, jwtManager, botMenu, api),
		List:         NewListService(jwtManager, bot, api),
	}
}
