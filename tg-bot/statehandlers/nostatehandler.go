package statehandlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
)

type NoStateHandler struct {
	Client       *redis.Client
	Bot          *tgbotapi.BotAPI
	StateHandler *StateHandler
}

// constants with response messages
const (
	startMessage = "Привет! Я бот для создания рецептов. Для регистрации введите /registration"
)

// HandleMessage функция для обработки команды в состоянии без состояния
func (nsh *NoStateHandler) HandleMessage(ctx context.Context, userID int64, command string) error {
	switch command {
	case startCommand:
		// Обработка команды /start
		err := nsh.Start(ctx, userID)
		if err != nil {
			return err
		}
	case registrationCommand:
		nsh.StateHandler.State = RegistrationState
		log.Println("registration")
	case createRecipeCommand:
		nsh.StateHandler.State = RecipeCreationState
	default:
		// Обработка неизвестной команды
		err := nsh.UnknownCommand(ctx, userID)
		if err != nil {
			return err
		}
	}
	err := nsh.StateHandler.setUserState(ctx, userID)
	if err != nil {
		return err
	}
	log.Println("state of :", nsh.StateHandler.State)
	return nil
}

// UnknownCommand функция для обработки неизвестной команды в состоянии без состояния
func (nsh *NoStateHandler) UnknownCommand(ctx context.Context, userID int64) error {
	// Вывод сообщения о том, что команда неизвестна
	_, err := nsh.Bot.Send(tgbotapi.NewMessage(userID, "Неизвестная команда"))
	if err != nil {
		return err
	}
	return nil

}

// Start функция для обработки команды /start в состоянии без состояния
func (nsh *NoStateHandler) Start(ctx context.Context, userID int64) error {
	// Вывод сообщения о том, что пользователь не зарегистрирован
	msg := tgbotapi.NewMessage(userID, startMessage)
	_, err := nsh.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// Cancel функция для обработки команды /cancel в состоянии без состояния
func (nsh *NoStateHandler) Cancel(ctx context.Context, userID int64) error {
	// Обработка команды /cancel
	// вывод сообщения о том что пользователь уже на начальном экране
	return nil
}

// HandleCallback функция для обработки callback-кнопок в состоянии без состояния
func (nsh *NoStateHandler) HandleCallback(ctx context.Context, userID int64, callbackData string) error {
	// Обработка callback-кнопок
	switch callbackData {
	case "cancel":
		// Обработка callback-кнопки "Отмена"
		log.Println("cancel")
	default:
		// Обработка неизвестной callback-кнопки
		return fmt.Errorf("unknown callback data")
	}
	return nil
}
