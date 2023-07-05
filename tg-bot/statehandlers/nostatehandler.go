package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
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
	StateHandler     *StateHandler
	UserStateManager interfaces.UserStateManager
	JwtManager       interfaces.JwtManager
	BotMenu          interfaces.BotMenu
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
		err := nsh.StateHandler.MealPlansList(ctx, userID)
		if err != nil {
			return state, err
		}
	case HelpCommand:
		err := nsh.Help(ctx, userID)
		if err != nil {
			return state, err
		}
	case RecipesListCommand:
		err := nsh.StateHandler.RecipesList(ctx, userID)
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

// LogOut функция для обработки команды /logout без состояния
func (nsh *NoStateHandler) LogOut(ctx context.Context, userID int64) error {
	//if user has already logged in - log out, else - send message that user is not logged in
	// Получение состояния пользователя
	token, err := nsh.JwtManager.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	if token != "" {
		// Удаление токена из базы данных
		err = nsh.JwtManager.DeleteUserJWTToken(ctx, userID)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(userID, "You logout successfully.")

		msg.ReplyMarkup = nsh.BotMenu.CreateMainMenu(ctx, userID)
		// Вывод сообщения о том, что пользователь вышел из аккаунта
		_, err := nsh.Bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	} else {
		// Вывод сообщения о том, что пользователь не зарегистрирован
		_, err := nsh.Bot.Send(tgbotapi.NewMessage(userID, "You are not logged in."))
		if err != nil {
			return err
		}
		return nil
	}

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

// Cancel функция для обработки команды /cancel в состоянии без состояния
func (nsh *NoStateHandler) Cancel(ctx context.Context, userID int64) error {
	// Обработка команды /cancel
	// вывод сообщения о том что пользователь уже на начальном экране
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
