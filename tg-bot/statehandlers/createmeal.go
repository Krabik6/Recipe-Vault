package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
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
	StateHandler     *StateHandler
	LocalState       createMealState
	UserStateManager interfaces.UserStateManager
	JwtManager       interfaces.JwtManager
	Client           *redis.Client
	Name             string
	Time             string
	Recipes          []int
}

const (
	createMealStateKey   = "create_meal_state:%d"
	createMealNameKey    = "create_meal_name:%d"
	createMealTimeKey    = "create_meal_time:%d"
	createMealRecipesKey = "create_meal_recipes:%d"
)

func (cms *CreateMealStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	fmt.Println(cms.LocalState, "Create meal state")
	query := update.CallbackQuery
	data := query.Data
	switch cms.LocalState {
	case CreateMealConfirmation:
		if data == "yes" {
			return cms.handleCreateMealConfirmYes(ctx, userID)
		} else {
			return cms.handleCreateMealConfirmNo(ctx, userID)
		}
	case CreateMealRecipes:
		if data == "confirm" {
			err := cms.handleCreateMealRecipes(userID)
			if err != nil {
				return err
			}
		} else if data == "yes" {
			return cms.handleCreateMealConfirmYes(ctx, userID)
		} else if data == "no" {
			return cms.handleCreateMealConfirmNo(ctx, userID)
		} else {
			err := cms.handleState(ctx, userID, update)
			log.Println("here")
			if err != nil {
				return err
			}

		}
	default:
		// msg to user: press one of the buttons or /cancel or print yes or no to confirm
		reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
		msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, reply)
		_, err := cms.Bot.Send(msg)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown state: damn%d", cms.LocalState)
	}

	//err := cms.handleState(ctx, userID)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
	return nil
}

// HandleMessage handles message for create meal state
func (cms *CreateMealStateHandler) HandleMessage(ctx context.Context, userID int64, update tgbotapi.Update) error {
	fmt.Println(cms.LocalState, "Create meal state")
	message := update.Message.Text
	if message == "/cancel" {
		return cms.handleCancel(ctx, userID)
	}
	switch cms.LocalState {
	//case NoCreateMealState:
	//	return cms.handleNoState(userID, update)
	case CreateMealName:
		cms.Name = message
	case CreateMealTime:
		cms.Time = message
	case CreateMealRecipes:
		if message == "confirm" {
			err := cms.handleCreateMealRecipes(userID)
			if err != nil {
				return err
			}
		} else {
			err := cms.handleState(ctx, userID, update)
			log.Println("here")
			if err != nil {
				return err
			}

		}

	case CreateMealConfirmation:
		if message == "yes" {
			return cms.handleCreateMealConfirmYes(ctx, userID)
		} else if message == "no" {
			return cms.handleCreateMealConfirmNo(ctx, userID)
		} else {
			// msg to user: press one of the buttons or /cancel or print yes or no to confirm
			reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
			msg := tgbotapi.NewMessage(update.Message.From.ID, reply)
			_, err := cms.Bot.Send(msg)
			if err != nil {
				return err
			}
			return fmt.Errorf("unknown state there: %d", cms.LocalState)
		}
	}

	err := cms.handleState(ctx, userID, update)
	if err != nil {
		return err
	}

	return nil
}

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
		msg.ReplyMarkup = cms.StateHandler.createCancelKeyboard()
		_, err := cms.StateHandler.Bot.Send(msg)
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

// handleCancel that cancel the creation of a recipe: delete the state from manager and send a message to the user
func (cms *CreateMealStateHandler) handleCancel(ctx context.Context, userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user

	cms.LocalState = NoCreateMealState
	err := cms.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = cms.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = cms.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Creation of meal plan canceled")
	msg.ReplyMarkup = cms.StateHandler.createMainMenu(ctx, userID)
	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// handleNoCreateMealState handles no create meal state
func (cms *CreateMealStateHandler) handleNoCreateMealState(userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	msg := tgbotapi.NewMessage(userID, "Enter the name of the meal plan")
	_, err := cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealName
	return nil
}

// handleCreateMealName handles create meal name state
func (cms *CreateMealStateHandler) handleCreateMealName(userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	//format 2022-05-21 00:0:31
	msg := tgbotapi.NewMessage(userID, "Enter the time of the meal plan in the format 2022-05-21 00:0:31")
	_, err := cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealTime
	return nil
}

// handleCreateMealTime handles create meal time state
func (cms *CreateMealStateHandler) handleCreateMealTime(ctx context.Context, userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	//format 2022-05-21 00:0:31
	err := cms.StateHandler.RecipesList(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Press the button with the recipe number to add it to the meal plan")

	row1 := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("confirm"),
		tgbotapi.NewKeyboardButton("/cancel"),
	}

	keyboard := tgbotapi.NewReplyKeyboard(row1)
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard

	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealRecipes
	return nil
}

// handleCreateMealRecipes handles create meal recipes state, send a message to the user with the the info about the meal plan
func (cms *CreateMealStateHandler) handleCreateMealRecipes(userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user
	reply := fmt.Sprintf("Meal plan: \n Name: %s \n Time: %s \n Recipes: %v", cms.Name, cms.Time, cms.Recipes)
	msg := tgbotapi.NewMessage(userID, reply)
	yesButton := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noButton := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yesButton, noButton))
	_, err := cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	cms.LocalState = CreateMealConfirmation
	return nil
}

// handleCreateMealConfirmYes handles the CreateMealConfirmation state when the user confirms the creation of the meal: set the state to NoCreateMealState and send a message to the user that the meal is created
func (cms *CreateMealStateHandler) handleCreateMealConfirmYes(ctx context.Context, userID int64) error {
	meal := &models.Meal{
		Name:    cms.Name,
		AtTime:  cms.Time,
		Recipes: cms.Recipes,
	}
	token, err := cms.JwtManager.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	client := &http.Client{}
	err = api.CreateMealPlan(client, *meal, token)
	if err != nil {
		return err
	}
	cms.LocalState = NoCreateMealState
	err = cms.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(userID, "Meal plan created on date: "+cms.Time)
	msg.ReplyMarkup = cms.StateHandler.createMainMenu(ctx, userID)
	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// handleCreateMealConfirmNo handles the CreateMealConfirmation state when the user does not confirm the creation of the meal: set the state to NoCreateMealState and send a message to the user that the meal is not created
func (cms *CreateMealStateHandler) handleCreateMealConfirmNo(ctx context.Context, userID int64) error {
	cms.LocalState = NoCreateMealState
	err := cms.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}

	err = cms.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(userID, "Meal plan not created")
	msg.ReplyMarkup = cms.StateHandler.createMainMenu(ctx, userID)
	_, err = cms.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
