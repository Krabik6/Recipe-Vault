package recipes

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
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

type CreateRecipeService struct {
	Bot *tgbotapi.BotAPI
	//StateHandler     *StateHandler
	State            *model.State
	LocalState       createRecipeState
	Client           *redis.Client
	UserStateManager interfaces.UserStateManager
	JwtManager       interfaces.JwtManager
	BotMenu          interfaces.BotMenu
	Api              *api.Api
	Title            string
	Description      string
	IsPublic         bool
	Cost             float64
	TimeToPrepare    int64
	Healthy          int
}

func NewCreateRecipeStateHandler(
	bot *tgbotapi.BotAPI,
	client *redis.Client,
	userManager interfaces.UserStateManager,
	jwtManager interfaces.JwtManager,
	botMenu interfaces.BotMenu,
	api *api.Api,

) *CreateRecipeService {
	if botMenu == nil {
		log.Fatal("botMenu is nil")
	}
	return &CreateRecipeService{
		Bot:              bot,
		Client:           client,
		UserStateManager: userManager,
		JwtManager:       jwtManager,
		BotMenu:          botMenu,
		Api:              api,
	}
}

// builder for title description isPublic cost timeToPrepare healthy that set values to struct fields (method)
func (crs *CreateRecipeService) BuildUserData(title string, description string, isPublic bool, cost float64, timeToPrepare int64, healthy int) {
	crs.Title = title
	crs.Description = description
	crs.IsPublic = isPublic
	crs.Cost = cost
	crs.TimeToPrepare = timeToPrepare
	crs.Healthy = healthy
}

// handleState handles the state
func (crs *CreateRecipeService) handleState(ctx context.Context, userID int64) error {
	//handle state, save data to manager
	switch crs.LocalState {
	case NoCreateRecipeState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		if crs.BotMenu == nil {
			log.Fatal("botMenu is nil")
		}
		msg.ReplyMarkup = crs.BotMenu.CreateCancelKeyboard() //todo change to interface
		_, err := crs.Bot.Send(msg)
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
