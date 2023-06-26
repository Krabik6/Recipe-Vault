package statehandlers

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

type RegistrationStateHandler struct {
	State        registrationState
	Client       *redis.Client
	StateHandler *StateHandler
	Bot          *tgbotapi.BotAPI
	Name         string
	Email        string
	Password     string
	ConfirmedPwd string
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

// constant for redis keys (user registration state, user registration data: name, email, password) like const userState = "user_state:%d"
const (
	userRegistrationState    = "user_registration_state:%d"
	userRegistrationName     = "user_registration_name:%d"
	userRegistrationEmail    = "user_registration_email:%d"
	userRegistrationPassword = "user_registration_password:%d"
)

// consts for buttons
const (
	confirmButton = "Confirm"
	cancelButton  = "Cancel"
)

// constants for callback queries data
const (
	confirmRegistrationCallbackData = "confirm_registration"
	cancelRegistrationCallbackData  = "cancel_registration"
)

// HandleCallbackQuery callback query handler
func (rs *RegistrationStateHandler) HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error {
	query := update.CallbackQuery
	data := query.Data
	log.Println("state: ", rs.State)
	switch rs.State {
	case RegistrationConfirmation:
		if data == confirmButton {
			log.Println("подтверждение регистрации")

			return rs.handleRegistrationComplete(ctx, userID, update)

		} else if data == cancelButton {
			message := tgbotapi.NewMessage(userID, "Регистрация отменена")
			_, err := rs.Bot.Send(message)
			if err != nil {
				return err
			}

			// Отмена регистрации, переход к начальному состоянию
			rs.State = NoRegistrationState
			err = rs.SetUserRegistrationState(ctx, userID)
			if err != nil {
				return err
			}
			err = rs.DeleteUserRegistrationData(ctx, userID)
			if err != nil {
				return err
			}
			rs.StateHandler.State = NoState
			err = rs.StateHandler.setUserState(ctx, userID)
			if err != nil {
				return err
			}

			return nil
		} else {
			return fmt.Errorf("unknown command")
		}
	default:
		// print query
		log.Println(query.Data)
		return fmt.Errorf("unknown registration state")
	}
}

func (rs *RegistrationStateHandler) HandleState(ctx context.Context, userID int64, update tgbotapi.Update) error {
	switch rs.State {
	case NoRegistrationState:
		msg := tgbotapi.NewMessage(userID, "To exit the current process, please press the \"Cancel\" button or enter \"/cancel\".") // Пустое текстовое сообщение
		msg.ReplyMarkup = rs.StateHandler.createCancelKeyboard()
		_, err := rs.StateHandler.Bot.Send(msg)
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
		rs.State = NoRegistrationState
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

func (rs *RegistrationStateHandler) HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error {
	//check if message is cancel then cancel registration
	if message == "/cancel" {

		rs.State = NoRegistrationState
		err := rs.SetUserRegistrationState(ctx, userID)
		if err != nil {
			return err
		}
		err = rs.DeleteUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		rs.StateHandler.State = NoState
		err = rs.StateHandler.setUserState(ctx, userID)
		if err != nil {
			return err
		}
		//add message to user
		msg := tgbotapi.NewMessage(userID, "Registration canceled")
		msg.ReplyMarkup = rs.StateHandler.createMainMenu(ctx, userID)
		_, err = rs.Bot.Send(msg)
		if err != nil {
			return err
		}

		log.Println("Registration canceled, state: ", rs.State)
		return nil
	}
	switch rs.State {
	case RegistrationName:
		rs.Name = message
	case RegistrationEmail:
		rs.Email = message
	case RegistrationPassword:
		rs.Password = message
	case RegistrationConfirmPassword:
		rs.ConfirmedPwd = message
	}

	// Обработка состояния
	err := rs.HandleState(ctx, userID, update)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return err
	}

	err = rs.SetUserRegistrationState(ctx, userID)
	if err != nil {
		return err
	}
	err = rs.SetUserRegistrationData(ctx, userID)
	if err != nil {
		return err
	}

	// Вывод текущего состояния
	fmt.Println("RegistrationStateHandler: ", rs.State)

	return nil
}

func (rs *RegistrationStateHandler) handleNoRegistrationState(userID int64) error {
	fmt.Println("Welcome! Please provide your name.")
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Welcome! Please provide your name."))
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на имя
	rs.State = RegistrationName
	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationName(userID int64) error {
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Hello, %s! Please provide your email.\n", rs.Name)))
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на email
	rs.State = RegistrationEmail

	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationEmail(ctx context.Context, userID int64) error {
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Username %s is registered. Please provide your password.\n", rs.Email)))
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на пароль
	rs.State = RegistrationPassword

	// Сохраняем состояние в Redis
	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationPassword(ctx context.Context, userID int64) error {
	// send message to user
	_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Password is set. Please confirm your password."))
	if err != nil {
		return err
	}

	// Меняем состояние регистрации на подтверждение пароля
	rs.State = RegistrationConfirmPassword

	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationConfirmPassword(ctx context.Context, userID int64, update tgbotapi.Update) error {
	if rs.ConfirmedPwd == rs.Password {
		_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Passwords match."))
		if err != nil {
			return err
		}
		rs.State = RegistrationConfirmation
		err = rs.handleRegistrationConfirmation(ctx, userID, update)
		if err != nil {
			return err
		}
	} else {
		_, err := rs.Bot.Send(tgbotapi.NewMessage(userID, "Password confirmation failed. Please try again. \n Provide your password."))
		if err != nil {
			return err
		}
		rs.State = RegistrationPassword
	}

	// Сохраняем состояние в Redis
	return nil
}

func (rs *RegistrationStateHandler) handleRegistrationConfirmation(ctx context.Context, userID int64, update tgbotapi.Update) error {
	reply := fmt.Sprintf("Please confirm your registration data.\nName: %s\nUsername: %s\nPassword: %s\n", rs.Name, rs.Email, rs.Password)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	// Кнопки подтверждения и отмены
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData(confirmButton, confirmButton)
	cancelBtn := tgbotapi.NewInlineKeyboardButtonData(cancelButton, cancelButton)
	// Добавляем кнопки в сообщение
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(confirmBtn, cancelBtn))
	// Отправляем сообщение пользователю
	_, err := rs.Bot.Send(msg)
	if err != nil {
		return err
	}
	// Меняем состояние регистрации на ожидание подтверждения
	rs.State = RegistrationConfirmation
	return nil

}

func (rs *RegistrationStateHandler) handleRegistrationComplete(ctx context.Context, userID int64, update tgbotapi.Update) error {

	user := model.SignUpCredentials{
		Username: rs.Email,
		Password: rs.Password,
		Name:     rs.Name,
	}
	client := &http.Client{}
	err := api.SignUp(rs.Bot, update, client, user)
	if err != nil {
		return err
	}
	// Обнуляем состояние регистрации
	rs.State = NoRegistrationState

	//delete state from redis
	err = rs.DeleteUserRegistrationState(ctx, userID)
	if err != nil {
		return err
	}

	err = rs.DeleteUserRegistrationData(ctx, userID)
	if err != nil {
		return err
	}

	err = rs.StateHandler.DeleteUserState(ctx, userID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(userID, "Registration complete!\n")
	msg.ReplyMarkup = rs.StateHandler.createMainMenu(ctx, userID)
	_, err = rs.Bot.Send(msg)

	if err != nil {
		return err
	}

	// Удаляем состояние из Redis, так как регистрация завершена
	return nil
}
