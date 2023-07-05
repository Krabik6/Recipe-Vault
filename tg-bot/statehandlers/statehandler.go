package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/manager"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"runtime/debug"
)

type StateHandler struct {
	//LocalState  LocalState
	Client           *redis.Client
	Bot              *tgbotapi.BotAPI
	UserStateManager interfaces.UserStateManager

	// Другие общие поля и методы, если необходимо
}

// UserData - struct for user data
//type UserData struct {
//	LocalState  model.LocalState
//	UserID int64
//	Update tgbotapi.Update
//}
//// NewUserData - constructor for UserData
//func NewUserData(state LocalState, userID int64, update tgbotapi.Update) *UserData {
//	return &UserData{
//		LocalState:  state,
//		UserID: userID,
//		Update: update,
//	}
//}

// NewStateHandler - constructor for StateHandler
func NewStateHandler(client *redis.Client, bot *tgbotapi.BotAPI) *StateHandler {
	return &StateHandler{
		Client:           client,
		Bot:              bot,
		UserStateManager: manager.NewRedisUserStateManager(client),
	}
}

// HandleMessage Метод для обработки входящих сообщений в соответствии с текущим состоянием
func (sh *StateHandler) HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error {
	//get user state from manager
	state, err := sh.UserStateManager.GetUserState(ctx, userID)
	if err != nil {
		return err
	}

	switch state {
	case model.NoState:
		noStateHandler := &NoStateHandler{
			Bot:              sh.Bot,
			StateHandler:     sh,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
		}
		state, err = noStateHandler.HandleMessage(ctx, userID, message, state)
		if err != nil {
			return err
		}
		if state != model.NoState {
			return sh.HandleMessage(ctx, userID, message, update)
		}
		return nil
	case model.RegistrationState:
		registrationHandler := &RegistrationStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
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
		registrationHandler.LocalState = regState
		return registrationHandler.HandleMessage(ctx, userID, message, update)
	case model.RecipeCreationState:
		recipeCreationHandler := &CreateRecipeStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
		}
		// set recipe creation data from manager
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
		// set recipe creation state from manager
		recipeCreationState, err := recipeCreationHandler.GetUserState(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.LocalState = recipeCreationState
		return recipeCreationHandler.HandleMessage(ctx, userID, update)
	case model.CreateMealState:
		createMealHandler := &CreateMealStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
		}
		// set recipe creation data from manager
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
		createMealHandler.LocalState = createMealStateGlobal
		log.Printf("Name: %s, Time: %s, Recipes: %v", name, time, recipes)
		return createMealHandler.HandleMessage(ctx, userID, update)

	case model.LogInState:
		logInHandler := &LoginStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
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

		logInHandler.LocalState = logInState
		return logInHandler.HandleMessage(ctx, userID, message, update)
	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}
}

// HandleCallbackQuery handle callback query
func (sh *StateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	state, err := sh.UserStateManager.GetUserState(ctx, userID)
	if err != nil {
		return err
	}
	switch state {
	case model.RecipeCreationState:
		recipeCreationHandler := &CreateRecipeStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
		}
		// set recipe creation data from manager
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
		// set recipe creation state from manager
		recipeCreationState, err := recipeCreationHandler.GetUserState(ctx, userID)
		if err != nil {
			return err
		}
		recipeCreationHandler.LocalState = recipeCreationState
		return recipeCreationHandler.HandleCallbackQuery(ctx, userID, update)
	case model.RegistrationState:
		rsh := RegistrationStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
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
		rsh.LocalState = regState
		err = rsh.HandleCallbackQuery(ctx, userID, update)
		if err != nil {
			return err
		}
	case model.LogInState:
		logInHandler := &LoginStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
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

		logInHandler.LocalState = logInState
		return logInHandler.HandleCallbackQuery(ctx, userID, update)
	case model.CreateMealState:
		createMealHandler := &CreateMealStateHandler{
			Client:           sh.Client,
			StateHandler:     sh,
			Bot:              sh.Bot,
			UserStateManager: manager.NewRedisUserStateManager(sh.Client),
		}
		// set recipe creation data from manager
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
		createMealHandler.LocalState = createMealStateGlobal
		log.Printf("Name: %s, Time: %s, Recipes: %v", name, time, recipes)
		return createMealHandler.HandleCallbackQuery(ctx, userID, update)

	default:
		//print stack trace
		debug.Stack()
		return fmt.Errorf("unknown state")
	}

	return nil
}
