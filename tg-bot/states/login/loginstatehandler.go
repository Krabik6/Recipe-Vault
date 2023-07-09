package login

//file name: loginstatehandler.go
import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"runtime/debug"
)

type LoginStateHandler struct {
	LocalState       loginState
	Client           *redis.Client
	Bot              *tgbotapi.BotAPI
	JwtManager       interfaces.JwtManager
	UserStateManager interfaces.UserStateManager
	BotMenu          interfaces.BotMenu
	Api              *api.Api
	Username         string
	Password         string
}

func NewLoginStateHandler(
	client *redis.Client,
	bot *tgbotapi.BotAPI,
	jwtManager interfaces.JwtManager,
	userStateManager interfaces.UserStateManager,
	botMenu interfaces.BotMenu,
	api *api.Api,
) *LoginStateHandler {
	return &LoginStateHandler{
		Client:           client,
		Bot:              bot,
		JwtManager:       jwtManager,
		UserStateManager: userStateManager,
		BotMenu:          botMenu,
		Api:              api,
	}
}

func (ls *LoginStateHandler) BuildUserData(userName string, password string) {
	ls.Username = userName
	ls.Password = password
}

type loginState int

const (
	NoLoginState loginState = iota
	LoginEmail
	LoginPassword
	LoginConfirmation
)

const (
	userLoginState    = "user_login_state:%d"
	userLoginEmail    = "user_login_username:%d"
	userLoginPassword = "user_login_password:%d"
)

func (ls *LoginStateHandler) HandleState(ctx context.Context, userID int64, update tgbotapi.Update) error {
	// state
	log.Println("state of handle state: ", ls.LocalState)
	switch ls.LocalState {
	case NoLoginState:
		err := ls.handleLoginEmail(userID)
		if err != nil {
			return err
		}
	case LoginEmail:
		err := ls.handleLoginPassword(userID)
		if err != nil {
			return err
		}
	case LoginPassword:
		err := ls.handleLoginConfirmation(ctx, userID, update)
		if err != nil {
			return err
		}

	case LoginConfirmation:
		msg := tgbotapi.NewMessage(userID, "Подтвердите вход в аккаунт")
		_, err := ls.Bot.Send(msg)
		if err != nil {
			return err
		}
	default:
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		// Save state to Redis
		ls.LocalState = NoLoginState
		err := ls.SetUserLoginState(ctx, userID)
		if err != nil {
			return err
		}
		err = ls.DeleteUserLoginData(ctx, userID)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown login state")
	}
	err := ls.SetUserLoginData(ctx, userID)
	if err != nil {
		return err
	}
	err = ls.SetUserLoginState(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
