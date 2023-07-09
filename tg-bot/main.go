package main

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/bot_buttons"
	"github.com/Krabik6/meal-schedule/tg-bot/manager"
	"github.com/Krabik6/meal-schedule/tg-bot/statehandlers"
	"github.com/Krabik6/meal-schedule/tg-bot/states/expectation"
	"github.com/Krabik6/meal-schedule/tg-bot/states/login"
	"github.com/Krabik6/meal-schedule/tg-bot/states/meals"
	"github.com/Krabik6/meal-schedule/tg-bot/states/recipes"
	"github.com/Krabik6/meal-schedule/tg-bot/states/registration"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	expiration = 7 * 24 * time.Hour
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
	ctx := context.Background()
	// Получение обновлений от Telegram API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	jwtManager := manager.NewJwtManager(redisClient, expiration)
	botMenu := bot_buttons.NewBotMenu(bot, jwtManager)
	botMenu.CreateMainMenu(ctx, 0)
	apis := api.NewApi()

	userStateManager := manager.NewUserStateManager(redisClient)
	if botMenu == nil {
		log.Fatal("botMenu is nil")
	}

	recipesService := recipes.NewRecipesService(bot, redisClient, userStateManager, jwtManager, botMenu, apis)
	mealServices := meals.NewMealsService(bot, userStateManager, jwtManager, redisClient, botMenu, recipesService, apis)
	loginService := login.NewLoginService(redisClient, bot, jwtManager, userStateManager, botMenu, apis)
	registerService := registration.NewRegistrationService(bot, redisClient, userStateManager, botMenu, apis)
	noStateService := expectation.NewNoStateHandler(bot, userStateManager, jwtManager, botMenu, mealServices, recipesService)

	sh := statehandlers.NewStateHandler(redisClient, bot, userStateManager, recipesService, mealServices, loginService, registerService, noStateService)
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
