package statehandlers

//file name: loginstatehandler.go
import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"runtime/debug"
)

type LoginStateHandler struct {
	State        loginState
	Client       *redis.Client
	StateHandler *StateHandler
	Bot          *tgbotapi.BotAPI
	Username     string
	Password     string
}

type loginState int

const (
	NoLoginState loginState = iota
	LoginEmail
	LoginPassword
	LoginConfirmation
)

// constant for redis keys (user login state, user login data: username, password) like const userState = "user_state:%d"

// consts for buttons
//const (
//	confirmButton = "Confirm"
//	cancelButton  = "Cancel"
//)

// constants for callback queries data
const (
	confirmLoginCallbackData = "confirm_login"
	cancelLoginCallbackData  = "cancel_login"
)

const (
	userLoginState    = "user_login_state:%d"
	userLoginEmail    = "user_login_username:%d"
	userLoginPassword = "user_login_password:%d"
)

// HandleCallbackQuery callback query handler
func (ls *LoginStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	query := update.CallbackQuery
	data := query.Data
	log.Println("state: ", ls.State)
	switch ls.State {
	case LoginConfirmation:
		if data == confirmButton {
			log.Println("Login confirmation")
			return ls.handleLoginComplete(ctx, userID, update)
		} else if data == cancelButton {
			message := tgbotapi.NewMessage(userID, "Вход в аккаунт отменен")
			_, err := ls.Bot.Send(message)
			if err != nil {
				return err
			}
			ls.State = NoLoginState
			err = ls.SetUserLoginState(ctx, userID)
			if err != nil {
				return err
			}
			err = ls.DeleteUserLoginData(ctx, userID)
			if err != nil {
				return err
			}
			ls.StateHandler.State = NoState
			err = ls.StateHandler.setUserState(ctx, userID)
			if err != nil {
				return err
			}
			return nil
		} else {

			return fmt.Errorf("unknown command")
		}
	case NoLoginState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = ls.StateHandler.createCancelKeyboard()
		_, err := ls.StateHandler.Bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	default:
		// print query
		log.Println(query.Data)
		return fmt.Errorf("unknown login state")
	}
}

func (ls *LoginStateHandler) HandleState(ctx context.Context, userID int64, update tgbotapi.Update) error {
	// state
	log.Println("state of handle state: ", ls.State)
	switch ls.State {
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
		ls.State = NoLoginState
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

func (ls *LoginStateHandler) HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error {
	// Check if the message is "/cancel" to cancel the login process
	if message == "/cancel" {
		ls.State = NoLoginState
		err := ls.SetUserLoginState(ctx, userID)
		if err != nil {
			return err
		}
		err = ls.DeleteUserLoginData(ctx, userID)
		if err != nil {
			return err
		}
		ls.StateHandler.State = NoState
		err = ls.StateHandler.setUserState(ctx, userID)
		if err != nil {
			return err
		}
		// Add message to the user
		msg := tgbotapi.NewMessage(userID, "Login canceled")
		msg.ReplyMarkup = ls.StateHandler.createMainMenu(ctx, userID)
		_, err = ls.Bot.Send(msg)
		if err != nil {
			return err
		}
		log.Println("Login canceled, state: ", ls.State)
		return nil
	}
	switch ls.State {
	case NoLoginState:
		//send cancel button without message
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = ls.StateHandler.createCancelKeyboard()
		_, err := ls.StateHandler.Bot.Send(msg)
		if err != nil {
			return err
		}
	case LoginEmail:
		ls.Username = message
	case LoginPassword:
		ls.Password = message
	}

	// Handle state
	err := ls.HandleState(ctx, userID, update)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return err
	}

	err = ls.SetUserLoginState(ctx, userID)
	if err != nil {
		return err
	}
	err = ls.SetUserLoginData(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ls *LoginStateHandler) handleNoLoginState(userID int64) error {
	// Send message to the user
	_, err := ls.Bot.Send(tgbotapi.NewMessage(userID, "Welcome! Please enter your username."))
	if err != nil {
		return err
	}
	// Change login state to username
	ls.State = LoginEmail
	return nil
}

func (ls *LoginStateHandler) handleLoginEmail(userID int64) error {
	// Send message to the user
	_, err := ls.Bot.Send(tgbotapi.NewMessage(userID, "Please enter your email."))
	if err != nil {
		return err
	}
	// Change login state to password
	ls.State = LoginEmail
	return nil
}

func (ls *LoginStateHandler) handleLoginPassword(userID int64) error {
	// Send message to the user
	_, err := ls.Bot.Send(tgbotapi.NewMessage(userID, "Email accepted. Please enter your password. "))
	if err != nil {
		return err
	}
	// Change login state to confirmation
	ls.State = LoginPassword
	return nil
}

// handleLoginConfirmation handles the login confirmation
func (ls *LoginStateHandler) handleLoginConfirmation(ctx context.Context, userID int64, update tgbotapi.Update) error {
	// Send message to the user like above, but with login data instead of registration data
	reply := fmt.Sprintf("Please confirm your login data.\nUsername: %s\nPassword: %s\n", ls.Username, ls.Password)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	// Confirm and cancel buttons
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData(confirmButton, confirmButton)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData(cancelButton, cancelButton)
	// Add buttons to the message
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(confirmBtn, cancelBtn))
	// Send message to the user
	_, err := ls.Bot.Send(msg)
	if err != nil {
		return err
	}
	// Change login state to confirmation
	ls.State = LoginConfirmation
	return nil
}

func (ls *LoginStateHandler) handleLoginComplete(ctx context.Context, userID int64, update tgbotapi.Update) error {
	user := model.LoginCredentials{
		Username: ls.Username,
		Password: ls.Password,
	}
	client := &http.Client{}
	token, err := api.Login(client, user)
	if err != nil {
		return fmt.Errorf("error logging in: %v", err)
	}
	err = ls.StateHandler.SetUserJWTToken(ctx, userID, token)
	if err != nil {
		return err
	}
	// Reset login state
	ls.State = NoLoginState

	// Delete state from Redis
	err = ls.DeleteUserLoginState(ctx, userID)
	if err != nil {
		return err
	}

	err = ls.DeleteUserLoginData(ctx, userID)
	if err != nil {
		return err
	}

	err = ls.StateHandler.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	// Add message to the user with the token
	msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Login complete! Your token is: %s\n", token))
	msg.ReplyMarkup = ls.StateHandler.createMainMenu(ctx, userID)

	_, err = ls.Bot.Send(msg)
	if err != nil {
		return err
	}
	// Remove state from Redis as login is completed
	return nil
}
