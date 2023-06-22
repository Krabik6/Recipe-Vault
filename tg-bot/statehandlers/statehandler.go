package statehandlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"runtime/debug"
)

type StateHandler struct {
	State  State
	Client *redis.Client
	Bot    *tgbotapi.BotAPI
	// Другие общие поля и методы, если необходимо
}

type State int

const (
	NoState State = iota
	RegistrationState
	RecipeCreationState
	// Другие состояния
)

// constants for commands (start, registration, etc)
const (
	startCommand        = "/start"
	registrationCommand = "/registration"
	createRecipeCommand = "/create_recipe"
	cancelCommand       = "/cancel"
)

// constant for redis key (user state)
const userState = "user_state:%d"

// Метод для обработки входящих сообщений в соответствии с текущим состоянием
func (sh *StateHandler) HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error {
	switch sh.State {
	case NoState:
		noStateHandler := &NoStateHandler{
			Client:       sh.Client,
			Bot:          sh.Bot,
			StateHandler: sh,
		}
		err := noStateHandler.HandleMessage(ctx, userID, message)
		if err != nil {
			return err
		}
		if sh.State != NoState {
			return sh.HandleMessage(ctx, userID, message, update)
		}
		return nil
	case RegistrationState:
		registrationHandler := &RegistrationStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		name, email, password, err := registrationHandler.GetUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		registrationHandler.Name = name
		registrationHandler.Email = email
		registrationHandler.Password = password
		regState, err := registrationHandler.GetUserRegistrationState(ctx, userID)
		if err != nil {
			return err
		}
		registrationHandler.State = regState
		log.Println(regState, "regstate")
		log.Println(name, email, password)
		return registrationHandler.HandleMessage(ctx, userID, message, update)
	case RecipeCreationState:
		recipeCreationHandler := &RecipeCreationStateHandler{
			Client: sh.Client,
		}
		return recipeCreationHandler.HandleMessage(ctx, userID, message)
	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}
}

// handle callback query
func (sh *StateHandler) HandleCallbackQuery(ctx context.Context, userID int64, query tgbotapi.CallbackQuery) error {
	switch sh.State {
	case RecipeCreationState:
		//todo: handle callback query in recipe creation state
	case RegistrationState:
		rsh := RegistrationStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		name, email, password, err := rsh.GetUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		rsh.Name = name
		rsh.Email = email
		rsh.Password = password
		regState, err := rsh.GetUserRegistrationState(ctx, userID)
		if err != nil {
			return err
		}
		rsh.State = regState
		log.Println(regState, "regstate")
		log.Println(name, email, password)

		err = rsh.HandleCallbackQuery(ctx, userID, query)
		if err != nil {
			return err
		}

	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}

	return nil
}

func HandleCommand(ctx context.Context, update tgbotapi.Update, redisClient *redis.Client, bot *tgbotapi.BotAPI) error {
	stateHandler := &StateHandler{
		Client: redisClient,
		Bot:    bot,
	}
	var userID int64

	// check update type
	if update.Message != nil {
		userID = update.Message.Chat.ID
		// Getting user state from redis
		_, err := stateHandler.getUserState(ctx, userID)
		if err != nil {
			return err
		}
		return stateHandler.HandleMessage(ctx, userID, update.Message.Text, update)
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.Message.Chat.ID
		// Getting user state from redis
		_, err := stateHandler.getUserState(ctx, userID)
		if err != nil {
			return err
		}
		return stateHandler.HandleCallbackQuery(ctx, userID, *update.CallbackQuery)
	}

	return nil
}
