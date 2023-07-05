package statehandlers

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

type createRecipeState int

const (
	NoCreateRecipeState createRecipeState = iota
	CreateRecipeTitle
	CreateRecipeDescription
	CreateRecipeIsPublic
	CreateRecipeCost
	CreateRecipeTimeToPrepare
	CreateRecipeHealthy
	CreateRecipeConfirmation
	CreateRecipeComplete
)

type CreateRecipeStateHandler struct {
	Bot              *tgbotapi.BotAPI
	StateHandler     *StateHandler
	State            *model.State
	LocalState       createRecipeState
	Client           *redis.Client
	UserStateManager interfaces.UserStateManager
	JwtManager       interfaces.JwtManager
	Title            string
	Description      string
	IsPublic         bool
	Cost             float64
	TimeToPrepare    int64
	Healthy          int
}

// consts for states of creating a recipe
/*
const (
	userLoginState    = "user_login_state:%d"
	userLoginEmail    = "user_login_username:%d"
	userLoginPassword = "user_login_password:%d"
)

like here
*/

const (
	createRecipeStateKey         = "create_recipe_state:%d"
	createRecipeTitleKey         = "create_recipe_title:%d"
	createRecipeDescriptionKey   = "create_recipe_description:%d"
	createRecipeIsPublicKey      = "create_recipe_isPublic:%d"
	createRecipeCostKey          = "create_recipe_cost:%d"
	createRecipeTimeToPrepareKey = "create_recipe_timeToPrepare:%d"
	createRecipeHealthyKey       = "create_recipe_healthy:%d"
)

// HandleMessage handles the message
func (crs *CreateRecipeStateHandler) HandleMessage(ctx context.Context, userID int64, update tgbotapi.Update) error {
	//print current state to console
	fmt.Printf("current state: %d\n", crs.LocalState)
	message := update.Message.Text
	if message == "/cancel" {
		return crs.handleCancel(ctx, userID)
	}
	switch crs.LocalState {
	case CreateRecipeTitle:
		crs.Title = message
	case CreateRecipeDescription:
		crs.Description = message
	case CreateRecipeIsPublic:
		if message == "yes" {
			crs.IsPublic = true
		} else {
			crs.IsPublic = false
		}
	case CreateRecipeCost:
		cost, err := strconv.ParseFloat(message, 64)
		if err != nil {
			return fmt.Errorf("error parsing cost: %v", err)
		}
		crs.Cost = cost
	case CreateRecipeTimeToPrepare:
		timeToPrepare, err := strconv.ParseInt(message, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing timeToPrepare: %v", err)
		}
		crs.TimeToPrepare = timeToPrepare
	case CreateRecipeHealthy:
		healthy, err := strconv.Atoi(message)
		if err != nil {
			return fmt.Errorf("error parsing healthy: %v", err)
		}
		crs.Healthy = healthy
	case CreateRecipeConfirmation:
		if message == "yes" {
			return crs.handleCreateRecipeConfirmYes(ctx, userID)
		} else if message == "no" {
			return crs.handleCreateRecipeConfirmNo(ctx, userID)
		} else {
			// msg to user: press one of the buttons or /cancel or print yes or no to confirm
			reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
			msg := tgbotapi.NewMessage(update.Message.From.ID, reply)
			_, err := crs.Bot.Send(msg)
			if err != nil {
				return err
			}
			return fmt.Errorf("unknown state: %d", crs.LocalState)
		}
	}

	//handle state
	err := crs.handleState(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// Handle CallbackQuery
func (crs *CreateRecipeStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	fmt.Printf("current state: %d\n", crs.LocalState)

	query := update.CallbackQuery
	data := query.Data
	switch crs.LocalState {
	case CreateRecipeIsPublic:
		if data == "yes" {
			crs.IsPublic = true
		} else {
			crs.IsPublic = false
		}
	case CreateRecipeHealthy:
		switch data {
		case "1":
			crs.Healthy = 1
		case "2":
			crs.Healthy = 2
		case "3":
			crs.Healthy = 3
		default:
			return fmt.Errorf("unknown healthy: %s", data)
		}
	case CreateRecipeConfirmation:
		if data == "yes" {
			return crs.handleCreateRecipeConfirmYes(ctx, userID)
		} else {
			return crs.handleCreateRecipeConfirmNo(ctx, userID)
		}
	default:
		// msg to user: press one of the buttons or /cancel or print yes or no to confirm
		reply := fmt.Sprintf("press one of the buttons or /cancel or print yes or no to confirm")
		msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, reply)
		_, err := crs.Bot.Send(msg)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown state: %d", crs.LocalState)
	}

	err := crs.handleState(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// handleState handles the state
func (crs *CreateRecipeStateHandler) handleState(ctx context.Context, userID int64) error {
	//handle state, save data to manager
	switch crs.LocalState {
	case NoCreateRecipeState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = crs.StateHandler.createCancelKeyboard()
		_, err := crs.StateHandler.Bot.Send(msg)
		if err != nil {
			return err
		}
		err = crs.handleNoCreateRecipeState(userID)
		if err != nil {
			return err
		}
	case CreateRecipeTitle:
		err := crs.handleCreateRecipeTitle(userID)
		if err != nil {
			return err
		}
	case CreateRecipeDescription:
		err := crs.handleCreateRecipeDescription(userID)
		if err != nil {
			return err
		}
	case CreateRecipeIsPublic:
		err := crs.handleCreateRecipeIsPublic(userID)
		if err != nil {
			return err
		}
	case CreateRecipeCost:
		err := crs.handleCreateRecipeCost(userID)
		if err != nil {
			return err
		}
	case CreateRecipeTimeToPrepare:
		err := crs.handleCreateRecipeTimeToPrepare(userID)
		if err != nil {
			return err
		}
	case CreateRecipeHealthy:
		err := crs.handleCreateRecipeHealthy(userID)
		if err != nil {
			return err
		}
	//case CreateRecipeConfirmation:
	//	err := crs.handleCreateRecipeConfirmation(ctx, userID)
	//	if err != nil {
	//		return err
	//	}
	default:
		return fmt.Errorf("unknown state: %d", crs.LocalState)
	}

	//save data to manager
	err := crs.SetUserData(ctx, userID)
	if err != nil {
		return err
	}

	//save state to manager
	err = crs.SetUserState(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// handleNoCreateRecipeState handles the NoCreateRecipeState state: set the state to CreateRecipeTitle and send a message to the user to ask for the title
func (crs *CreateRecipeStateHandler) handleNoCreateRecipeState(userID int64) error {
	// set the state to CreateRecipeTitle and send a message to the user to ask for the title

	msg := tgbotapi.NewMessage(userID, "What is the title of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeTitle

	return nil
}

// handleCreateRecipeTitle handles the CreateRecipeTitle state: set the state to CreateRecipeDescription and send a message to the user to ask for the description
func (crs *CreateRecipeStateHandler) handleCreateRecipeTitle(userID int64) error {
	// set the state to CreateRecipeDescription and send a message to the user to ask for the description

	msg := tgbotapi.NewMessage(userID, "What is the description of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeDescription

	return nil
}

// handleCreateRecipeDescription handles the CreateRecipeDescription state: set the state to CreateRecipeIsPublic and send a message to the user to ask if the recipe is public
func (crs *CreateRecipeStateHandler) handleCreateRecipeDescription(userID int64) error {
	// set the state to CreateRecipeIsPublic and send a message with the keyboard to the user to ask if the recipe is public
	reply := fmt.Sprintf("Is the recipe public?")
	msg := tgbotapi.NewMessage(userID, reply)
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData("no", "no")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(confirmBtn, cancelBtn))
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeIsPublic

	return nil
}

// handleCreateRecipeIsPublic handles the CreateRecipeIsPublic state: set the state to CreateRecipeCost and send a message to the user to ask for the cost
func (crs *CreateRecipeStateHandler) handleCreateRecipeIsPublic(userID int64) error {
	// set the state to CreateRecipeCost and send a message to the user to ask for the cost

	msg := tgbotapi.NewMessage(userID, "What is the cost of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeCost

	return nil
}

// handleCreateRecipeCost handles the CreateRecipeCost state: set the state to CreateRecipeTimeToPrepare and send a message to the user to ask for the time to prepare
func (crs *CreateRecipeStateHandler) handleCreateRecipeCost(userID int64) error {
	// set the state to CreateRecipeTimeToPrepare and send a message to the user to ask for the time to prepare

	msg := tgbotapi.NewMessage(userID, "What is the time to prepare of the recipe?")
	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}
	crs.LocalState = CreateRecipeTimeToPrepare

	return nil
}

// handleCreateRecipeTimeToPrepare handles the CreateRecipeTimeToPrepare state: set the state to CreateRecipeHealthy and send a message to the user to ask if the recipe is healthy
func (crs *CreateRecipeStateHandler) handleCreateRecipeTimeToPrepare(userID int64) error {
	// set the state to CreateRecipeHealthy and send a message with the 3 buttons to the user to ask how healthy is the recipe (not healthy, healthy, very healthy)
	reply := fmt.Sprintf("How healthy is the recipe?")
	msg := tgbotapi.NewMessage(userID, reply)
	notHealthyBtn := tgbotapi.NewInlineKeyboardButtonData("not healthy", "1")
	healthyBtn := tgbotapi.NewInlineKeyboardButtonData("healthy", "2")
	veryHealthyBtn := tgbotapi.NewInlineKeyboardButtonData("very healthy", "3")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(notHealthyBtn, healthyBtn, veryHealthyBtn))

	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	crs.LocalState = CreateRecipeHealthy

	return nil
}

// handleCreateRecipeHealthy handles the CreateRecipeHealthy state: set the state to CreateRecipeConfirmation and send a message to the user to ask for the confirmation
func (crs *CreateRecipeStateHandler) handleCreateRecipeHealthy(userID int64) error {
	// set the state to CreateRecipeConfirmation and send a message with info about recipe and with the 2 buttons to the user to ask for the confirmation (yes, no)
	reply := fmt.Sprintf("Recipe info:\nTitle: %s\nDescription: %s\nIs public: %t\nCost: %d\nTime to prepare: %d\nHealthy: %s\n\nIs the info correct?", crs.Title, crs.Description, crs.IsPublic, crs.Cost, crs.TimeToPrepare, crs.Healthy)
	// create repy with beutiful format of the recipe info
	reply = fmt.Sprintf("Recipe info:\nTitle: %s\nDescription: %s\nIs public: %t\nCost: %.2f\nTime to prepare: %d\nHealthy: %d\n\nIs the info correct?", crs.Title, crs.Description, crs.IsPublic, crs.Cost, crs.TimeToPrepare, crs.Healthy)
	msg := tgbotapi.NewMessage(userID, reply)
	yesBtn := tgbotapi.NewInlineKeyboardButtonData("yes", "yes")
	noBtn := tgbotapi.NewInlineKeyboardButtonData("no", "no")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yesBtn, noBtn))

	_, err := crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	crs.LocalState = CreateRecipeConfirmation

	return nil
}

// handleCreateRecipeConfirmYes handles the CreateRecipeConfirmation state when the user confirms the creation of the recipe: set the state to NoCreateRecipeState and send a message to the user that the recipe is created
func (crs *CreateRecipeStateHandler) handleCreateRecipeConfirmYes(ctx context.Context, userID int64) error {
	// set the state to NoCreateRecipeState and send a message to the user that the recipe is created

	// create the recipe
	recipe := model.CreateRecipeInput{
		Title:         crs.Title,
		Description:   crs.Description,
		IsPublic:      crs.IsPublic,
		Cost:          crs.Cost,
		TimeToPrepare: crs.TimeToPrepare,
		Healthy:       crs.Healthy,
	}
	client := &http.Client{}
	//gwt jwt token
	token, err := crs.JwtManager.GetUserJWTToken(ctx, userID)

	recipeId, err := api.CreateRecipe(client, recipe, token)
	if err != nil {
		return err
	}

	crs.LocalState = NoCreateRecipeState
	err = crs.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Recipe created with id %d", recipeId))
	msg.ReplyMarkup = crs.StateHandler.createMainMenu(ctx, userID)
	_, err = crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// handleCreateRecipeConfirmNo handles the CreateRecipeConfirmation state when the user doesn't confirm the creation of the recipe: set the state to NoCreateRecipeState and send a message to the user that the recipe is not created also delete
func (crs *CreateRecipeStateHandler) handleCreateRecipeConfirmNo(ctx context.Context, userID int64) error {
	// set the state to NoCreateRecipeState and send a message to the user that the recipe is not created

	crs.LocalState = NoCreateRecipeState
	err := crs.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Creation of recipe canceled")
	_, err = crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// handleCancel that cancel the creation of a recipe: delete the state from manager and send a message to the user
func (crs *CreateRecipeStateHandler) handleCancel(ctx context.Context, userID int64) error {
	// local state to no state and delete the state from manager, send a message to the user

	crs.LocalState = NoCreateRecipeState
	err := crs.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.DeleteUserData(ctx, userID)
	if err != nil {
		return err
	}
	err = crs.UserStateManager.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Creation of recipe canceled")
	msg.ReplyMarkup = crs.StateHandler.createMainMenu(ctx, userID)
	_, err = crs.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
