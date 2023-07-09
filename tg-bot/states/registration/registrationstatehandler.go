package registration

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

type RegistrationStateHandler struct {
	LocalState       registrationState
	Client           *redis.Client
	UserStateManager interfaces.UserStateManager
	Bot              *tgbotapi.BotAPI
	BotMenu          interfaces.BotMenu
	Api              *api.Api
	Name             string
	Email            string
	Password         string
	ConfirmedPwd     string
}

func NewRegistrationStateHandler(
	client *redis.Client,
	userStateManager interfaces.UserStateManager,
	bot *tgbotapi.BotAPI,
	botMenu interfaces.BotMenu,
	api *api.Api,
) *RegistrationStateHandler {
	return &RegistrationStateHandler{
		Client:           client,
		UserStateManager: userStateManager,
		Bot:              bot,
		BotMenu:          botMenu,
		Api:              api,
	}
}

// builder for Name, email, password, confirmedPwd in method form
func (rs *RegistrationStateHandler) BuildUserData(name, email, password string) {
	rs.Name = name
	rs.Email = email
	rs.Password = password
}

type registrationState int

const (
	NoRegistrationState registrationState = iota
	RegistrationName
	RegistrationEmail
	RegistrationPassword
	RegistrationConfirmPassword
	// state for registration confirmation
	RegistrationConfirmation
)

// constant for manager keys (user registration state, user registration data: name, email, password) like const userState = "user_state:%d"
const (
	userRegistrationState    = "user_registration_state:%d"
	userRegistrationName     = "user_registration_name:%d"
	userRegistrationEmail    = "user_registration_email:%d"
	userRegistrationPassword = "user_registration_password:%d"
)

func (rs *RegistrationStateHandler) HandleState(ctx context.Context, userID int64, update tgbotapi.Update) error {
	switch rs.LocalState {
	case NoRegistrationState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = rs.BotMenu.CreateCancelKeyboard()
		_, err := rs.Bot.Send(msg)
		if err != nil {
			return err
		}
		err = rs.handleNoRegistrationState(userID)
		if err != nil {
			return err
		}
	case RegistrationName:
		err := rs.handleRegistrationName(userID)
		if err != nil {
			return err
		}
	case RegistrationEmail:
		err := rs.handleRegistrationEmail(ctx, userID)
		if err != nil {
			return err
		}
	case RegistrationPassword:
		err := rs.handleRegistrationPassword(ctx, userID)
		if err != nil {
			return err
		}
	case RegistrationConfirmPassword:
		err := rs.handleRegistrationConfirmPassword(ctx, userID, update)
		if err != nil {
			return err
		}
	case RegistrationConfirmation:
		// send message that user should confirm registration pressing button or sending /cancel
		msg := tgbotapi.NewMessage(userID, "Подтвердите регистрацию")
		_, err := rs.Bot.Send(msg)
		if err != nil {
			return err
		}
	default:
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		// Сохраняем состояние в Redis
		rs.LocalState = NoRegistrationState
		err := rs.SetUserRegistrationState(ctx, userID)
		if err != nil {
			return err
		}
		err = rs.DeleteUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown registration state")
	}
	err := rs.SetUserRegistrationData(ctx, userID)
	if err != nil {
		return err
	}
	err = rs.SetUserRegistrationState(ctx, userID)
	if err != nil {
		return err
	}
	log.Println("дошло до 215 строки")
	return nil

}
