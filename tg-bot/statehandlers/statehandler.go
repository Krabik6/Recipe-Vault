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
	State      State
	Registered bool
	Client     *redis.Client
	Bot        *tgbotapi.BotAPI
	// Другие общие поля и методы, если необходимо
}

type State int

const (
	NoState State = iota
	RegistrationState
	RecipeCreationState
	LogInState
	CreateMealState
	// Другие состояния
)

// constants for commands (start, registration, etc)

// constant for redis key (user state)
const userState = "user_state:%d"

// HandleMessage Метод для обработки входящих сообщений в соответствии с текущим состоянием
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
		return registrationHandler.HandleMessage(ctx, userID, message, update)
	case RecipeCreationState:
		recipeCreationHandler := &CreateRecipeStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		// set recipe creation data from redis
		title, description, isPublic, cost, timeToPrepare, healthy, err := recipeCreationHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.Title = title
		recipeCreationHandler.Description = description
		recipeCreationHandler.IsPublic = isPublic
		recipeCreationHandler.Cost = cost
		recipeCreationHandler.TimeToPrepare = timeToPrepare
		recipeCreationHandler.Healthy = healthy
		//Print
		log.Printf("Title: %s, Description: %s, IsPublic: %t, Cost: %d, TimeToPrepare: %d, Healthy: %t", title, description, isPublic, cost, timeToPrepare, healthy)
		// set recipe creation state from redis
		recipeCreationState, err := recipeCreationHandler.GetUserState(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.State = recipeCreationState
		return recipeCreationHandler.HandleMessage(ctx, userID, update)
	case CreateMealState:
		createMealHandler := &CreateMealStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		// set recipe creation data from redis
		name, time, recipes, err := createMealHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}

		createMealHandler.Name = name
		createMealHandler.Time = time
		createMealHandler.Recipes = recipes
		createMealStateGlobal, err := createMealHandler.GetUserState(ctx, userID)
		if err != nil {
			return err
		}
		createMealHandler.State = createMealStateGlobal
		log.Printf("Name: %s, Time: %s, Recipes: %v", name, time, recipes)
		return createMealHandler.HandleMessage(ctx, userID, update)

	case LogInState:
		logInHandler := &LoginStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		username, password, err := logInHandler.GetUserLoginData(ctx, userID)
		if err != nil {
			return err
		}

		logInHandler.Username = username
		logInHandler.Password = password
		logInState, err := logInHandler.GetUserLoginState(ctx, userID)
		if err != nil {
			return err
		}

		logInHandler.State = logInState
		return logInHandler.HandleMessage(ctx, userID, message, update)
	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}
}

// HandleCallbackQuery handle callback query
func (sh *StateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	log.Println("HandleCallbackQuery global state: ", sh.State)
	switch sh.State {
	case RecipeCreationState:
		recipeCreationHandler := &CreateRecipeStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		// set recipe creation data from redis
		title, description, isPublic, cost, timeToPrepare, healthy, err := recipeCreationHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.Title = title
		recipeCreationHandler.Description = description
		recipeCreationHandler.IsPublic = isPublic
		recipeCreationHandler.Cost = cost
		recipeCreationHandler.TimeToPrepare = timeToPrepare
		recipeCreationHandler.Healthy = healthy
		//Print
		log.Printf("Title: %s, Description: %s, IsPublic: %t, Cost: %d, TimeToPrepare: %d, Healthy: %t", title, description, isPublic, cost, timeToPrepare, healthy)
		// set recipe creation state from redis
		recipeCreationState, err := recipeCreationHandler.GetUserState(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.State = recipeCreationState
		return recipeCreationHandler.HandleCallbackQuery(ctx, userID, update)
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
		err = rsh.HandleCallbackQuery(ctx, userID, update)
		if err != nil {
			return err
		}
	case LogInState:
		logInHandler := &LoginStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		username, password, err := logInHandler.GetUserLoginData(ctx, userID)
		if err != nil {
			return err
		}

		logInHandler.Username = username
		logInHandler.Password = password
		logInState, err := logInHandler.GetUserLoginState(ctx, userID)
		if err != nil {
			return err
		}

		logInHandler.State = logInState
		return logInHandler.HandleCallbackQuery(ctx, userID, update)
	case CreateMealState:
		createMealHandler := &CreateMealStateHandler{
			Client:       sh.Client,
			StateHandler: sh,
			Bot:          sh.Bot,
		}
		// set recipe creation data from redis
		name, time, recipes, err := createMealHandler.GetUserData(ctx, userID)
		if err != nil {
			return err
		}

		createMealHandler.Name = name
		createMealHandler.Time = time
		createMealHandler.Recipes = recipes
		createMealStateGlobal, err := createMealHandler.GetUserState(ctx, userID)
		if err != nil {
			return err
		}
		createMealHandler.State = createMealStateGlobal
		log.Printf("Name: %s, Time: %s, Recipes: %v", name, time, recipes)
		return createMealHandler.HandleCallbackQuery(ctx, userID, update)

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
		state, err := stateHandler.getUserState(ctx, userID)
		if err != nil {
			return err
		}
		log.Println("global state: ", state)
		//create cancel button
		if err != nil {
			return err
		}

		return stateHandler.HandleMessage(ctx, userID, update.Message.Text, update)
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.Message.Chat.ID
		// Getting user state from redis
		state, err := stateHandler.getUserState(ctx, userID)
		if err != nil {
			return err
		}
		log.Println("global state: ", state)

		return stateHandler.HandleCallbackQuery(ctx, userID, update)
	}

	return nil
}
