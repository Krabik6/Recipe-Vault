package registration

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
)

/*


 */

func (rs *RegistrationStateHandler) HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error {
	//get state
	state, err := rs.GetUserRegistrationState(ctx, userID)
	if err != nil {
		return err
	}
	rs.LocalState = state
	//check if message is cancel then cancel registration
	if message == "/cancel" {

		rs.LocalState = NoRegistrationState
		err := rs.SetUserRegistrationState(ctx, userID)
		if err != nil {
			return err
		}
		err = rs.DeleteUserRegistrationData(ctx, userID)
		if err != nil {
			return err
		}
		err = rs.UserStateManager.DeleteUserState(ctx, userID)
		if err != nil {
			return err
		}
		//add message to user
		msg := tgbotapi.NewMessage(userID, "Registration canceled")
		msg.ReplyMarkup = rs.BotMenu.CreateMainMenu(ctx, userID)
		_, err = rs.Bot.Send(msg)
		if err != nil {
			return err
		}

		log.Println("Registration canceled, state: ", rs.LocalState)
		return nil
	}
	switch rs.LocalState {
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
	err = rs.HandleState(ctx, userID, update)
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
	fmt.Println("RegistrationStateHandler: ", rs.LocalState)

	return nil
}
