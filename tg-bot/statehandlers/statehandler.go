package statehandlers

import (
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/states/expectation"
	"github.com/Krabik6/meal-schedule/tg-bot/states/login"
	"github.com/Krabik6/meal-schedule/tg-bot/states/meals"
	"github.com/Krabik6/meal-schedule/tg-bot/states/recipes"
	"github.com/Krabik6/meal-schedule/tg-bot/states/registration"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type StateHandler struct {
	//LocalState  LocalState
	Client           *redis.Client
	Bot              *tgbotapi.BotAPI
	UserStateManager interfaces.UserStateManager
	BotMenu          interfaces.BotMenu
	JwtManager       interfaces.JwtManager
	Recipes          *recipes.RecipesService
	Meals            *meals.MealsService
	Login            *login.LoginService
	Registration     *registration.RegistrationService
	NoState          *expectation.NoStateHandler
}

// NewStateHandler - constructor for StateHandler
func NewStateHandler(
	client *redis.Client,
	bot *tgbotapi.BotAPI,
	userStateManager interfaces.UserStateManager,
	recipesService *recipes.RecipesService,
	meals *meals.MealsService,
	login *login.LoginService,
	registration *registration.RegistrationService,
	noState *expectation.NoStateHandler,
) *StateHandler {
	return &StateHandler{
		Client:           client,
		Bot:              bot,
		UserStateManager: userStateManager,
		Recipes:          recipesService,
		Meals:            meals,
		Login:            login,
		Registration:     registration,
		NoState:          noState,
	}
}
