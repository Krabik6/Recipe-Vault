package statehandlers

//
//import (
//	"context"
//	"fmt"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/redis/go-redis/v9"
//	"runtime/debug"
//)
//
//type CancelButton struct {
//	Bot    *tgbotapi.BotAPI
//	Client *redis.Client
//}
//
//func (cb *redis.Client) HandleCancel(ctx context.Context, userID int64) error {
//	// Удаление состояния регистрации из Redis
//	err := cb.Client.Del(ctx, fmt.Sprintf("state:%d", userID)).Err()
//	if err != nil {
//		stackTrace := debug.Stack()
//		fmt.Println("Stack trace:", string(stackTrace))
//		return err
//	}
//
//	// Удаление данных пользователя из Redis
//	err = cb.Client.Del(ctx, fmt.Sprintf("user:%d", userID)).Err()
//	if err != nil {
//		stackTrace := debug.Stack()
//		fmt.Println("Stack trace:", string(stackTrace))
//		return err
//	}
//
//	// Вывод сообщения о успешной отмене
//	msg := tgbotapi.NewMessage(userID, "Регистрация отменена.")
//	_, err = cb.Bot.Send(msg)
//	if err != nil {
//		stackTrace := debug.Stack()
//		fmt.Println("Stack trace:", string(stackTrace))
//		return err
//	}
//
//	return nil
//}
