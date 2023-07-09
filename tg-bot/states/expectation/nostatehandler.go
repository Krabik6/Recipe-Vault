package expectation

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"github.com/Krabik6/meal-schedule/tg-bot/states/meals"
	"github.com/Krabik6/meal-schedule/tg-bot/states/recipes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	StartMessage = "Привет! Я бот для создания рецептов.\n Список комманд: \n /start - начать работу с ботом \n /registration - зарегистрироваться \n /login - войти в аккаунт \n /create_recipe - создать рецепт  \n /logout - выйти из аккаунта"
)

const (
	HelpCommand         = "/help"
	RegistrationCommand = "/signup"
	CreateRecipeCommand = "/create_recipe"
	LogInCommand        = "/login"
	LogOutCommand       = "/logout"
	StartCommand        = "/start"
	CancelCommand       = "/cancel"
	RecipesListCommand  = "/recipes_list"
	CreateMealCommand   = "/create_meal"
	MealsListCommand    = "/meals_list"
)

type NoStateHandler struct {
	Bot              *tgbotapi.BotAPI
	UserStateManager interfaces.UserStateManager
	JwtManager       interfaces.JwtManager
	BotMenu          interfaces.BotMenu
	Meal             *meals.MealsService
	Recipe           *recipes.RecipesService
}

func NewNoStateHandler(
	bot *tgbotapi.BotAPI,
	userStateManager interfaces.UserStateManager,
	jwtManager interfaces.JwtManager,
	botMenu interfaces.BotMenu,
	meal *meals.MealsService,
	recipe *recipes.RecipesService,
) *NoStateHandler {
	return &NoStateHandler{
		Bot:              bot,
		UserStateManager: userStateManager,
		JwtManager:       jwtManager,
		BotMenu:          botMenu,
		Meal:             meal,
		Recipe:           recipe,
	}
}

// HandleMessage функция для обработки команды в состоянии без состояния
func (nsh *NoStateHandler) HandleMessage(ctx context.Context, userID int64, command string, state model.State) (model.State, error) {
	switch command {
	case StartCommand:
		// Вывод сообщения о том, что пользователь уже зарегистрирован
		err := nsh.Start(ctx, userID)
		if err != nil {
			return state, err
		}
	case MealsListCommand:
		err := nsh.Meal.MealPlansList(ctx, userID)
		if err != nil {
			return state, err
		}
	case HelpCommand:
		err := nsh.Help(ctx, userID)
		if err != nil {
			return state, err
		}
	case RecipesListCommand:
		err := nsh.Recipe.RecipesList(ctx, userID)
		if err != nil {
			return state, err
		}
	case LogOutCommand:
		err := nsh.LogOut(ctx, userID)
		if err != nil {
			return state, err
		}
	case CreateMealCommand:
		state = model.CreateMealState
	case RegistrationCommand:
		state = model.RegistrationState
	case CreateRecipeCommand:
		state = model.RecipeCreationState
	case LogInCommand:
		state = model.LogInState
	default:
		// Обработка неизвестной команды
		err := nsh.UnknownCommand(ctx, userID)
		if err != nil {
			return state, err
		}
	}
	err := nsh.UserStateManager.SetUserState(ctx, userID, state)
	if err != nil {
		return state, err
	}
	return state, nil
}

// UnknownCommand функция для обработки неизвестной команды в состоянии без состояния
func (nsh *NoStateHandler) UnknownCommand(ctx context.Context, userID int64) error {
	// Вывод сообщения о том, что команда неизвестна
	_, err := nsh.Bot.Send(tgbotapi.NewMessage(userID, "Неизвестная команда"))
	if err != nil {
		return err
	}
	return nil

}

// Start функция для обработки команды /start в состоянии без состояния, что будет отображаться при входе в бота и также отображает кнопки на боте
func (nsh *NoStateHandler) Start(ctx context.Context, userID int64) error {
	msg := tgbotapi.NewMessage(userID, StartMessage)

	msg.ReplyMarkup = nsh.BotMenu.CreateMainMenu(ctx, userID)
	_, err := nsh.Bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// Help функция для обработки команды /start в состоянии без состояния
func (nsh *NoStateHandler) Help(ctx context.Context, userID int64) error {
	// Вывод сообщения о том, что пользователь не зарегистрирован
	msg := tgbotapi.NewMessage(userID, StartMessage)
	_, err := nsh.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// HandleCallback функция для обработки callback-кнопок в состоянии без состояния
func (nsh *NoStateHandler) HandleCallback(ctx context.Context, userID int64, callbackData string) error {
	// Обработка callback-кнопок
	switch callbackData {
	case "cancel":
		// Обработка callback-кнопки "Отмена"
		log.Println("cancel")
	default:
		// Обработка неизвестной callback-кнопки
		return fmt.Errorf("unknown callback data")
	}
	return nil
}
