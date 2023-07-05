package main

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/statehandlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	// Команда для начала регистрации
	StartRegistrationCommand = "/registration"
	// Команда для отмены регистрации
	CancelRegistrationCommand = "/cancel"
)

func main() {
	// Создание клиента Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Пароль Redis, если применимо
		DB:       0,  // Номер базы данных Redis, если применимо
	})

	//check manager connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	//g
	botToken := "5790667960:AAHC1XU-IWF6aQ1p57fLdCu30_WXgys3vXo"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	//

	// Создание контекста
	ctx := context.Background()

	// Получение обновлений от Telegram API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	sh := statehandlers.NewStateHandler(redisClient, bot)
	log.Println("Bot is running...")

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {
			// Handle message updates
			err := sh.HandleMessage(ctx, update.Message.Chat.ID, update.Message.Text, update)
			if err != nil {
				// В случае ошибки отправляем сообщение пользователю и логируем ошибку
				message := fmt.Sprintf("Произошла ошибка: %s\n Пожалуйста попробуйте ещё раз.", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
				}
				log.Println("Error handling message:", err)
			}
		} else if update.CallbackQuery != nil {
			// Handle callback query updates
			err := sh.HandleCallbackQuery(ctx, update.CallbackQuery.Message.Chat.ID, update)
			if err != nil {
				// В случае ошибки отправляем сообщение пользователю и логируем ошибку
				message := fmt.Sprintf("Произошла ошибка: %s\n Пожалуйста попробуйте ещё раз.", err)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
				}
				log.Println("Error handling callback query:", err)
			}
		}

	}
}
