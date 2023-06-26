package main

import (
	"context"
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

	//check redis connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	botToken := "5790667960:AAHC1XU-IWF6aQ1p57fLdCu30_WXgys3vXo"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Создание контекста
	ctx := context.Background()

	// Получение обновлений от Telegram API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		//log.Println(update.CallbackQuery.Data)

		//update.CallbackQuery.Data
		// Обработка команд пользователя

		//command := update.Message.Text
		err := statehandlers.HandleCommand(ctx, update, redisClient, bot)
		if err != nil {
			spoilerErrorMessage := "Произошла ошибка: " + err.Error()

			msg := tgbotapi.NewMessage(update.FromChat().ID, spoilerErrorMessage)
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		// Вывод текущего состояния
	}
}
