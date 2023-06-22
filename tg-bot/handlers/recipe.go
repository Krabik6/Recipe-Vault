package handlers

//
//import (
//	"fmt"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/redis/go-redis/v9"
//	"log"
//	"strings"
//)
//
//func handleRecipeIdleState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	if update.Message.Text == "/createRecipe" {
//		reply := "Добро пожаловать в процесс создания рецепта!\nПожалуйста, введите название рецепта:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//
//		recipe.RecipeState = StateWaitingForRecipeName
//		err := saveRecipeToRedis(redisClient, recipe)
//		if err != nil {
//			log.Println("Ошибка при сохранении рецепта в Redis:", err)
//		}
//	} else if update.Message.Text == "/cancel" {
//		handleCancelState(bot, redisClient, user, recipe, update)
//	}
//}
//
//func handleRecipeNameState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	recipeName := strings.TrimSpace(update.Message.Text)
//	if recipeName == "" {
//		reply := "Название рецепта не может быть пустым. Пожалуйста, введите название рецепта:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//		return
//	}
//
//	recipe.RecipeName = recipeName
//
//	reply := "Отлично! Теперь введите описание рецепта:"
//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//	bot.Send(msg)
//
//	recipe.RecipeState = StateWaitingForRecipeDescription
//	err := saveRecipeToRedis(redisClient, recipe)
//	if err != nil {
//		log.Println("Ошибка при сохранении рецепта в Redis:", err)
//	}
//}
//
//func handleRecipeDescriptionState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	description := strings.TrimSpace(update.Message.Text)
//	if description == "" {
//		reply := "Описание рецепта не может быть пустым. Пожалуйста, введите описание рецепта:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//		return
//	}
//
//	recipe.Description = description
//
//	reply := fmt.Sprintf("Вы ввели следующие данные для рецепта:\nНазвание: %s\nОписание: %s\n\nПожалуйста, подтвердите, что все верно.",
//		recipe.RecipeName, recipe.Description)
//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("Да", "recipe_confirm"),
//			tgbotapi.NewInlineKeyboardButtonData("Нет", "recipe_restart"),
//		),
//	)
//	bot.Send(msg)
//
//	recipe.RecipeState = StateRecipeConfirmation
//	err := saveRecipeToRedis(redisClient, recipe)
//	if err != nil {
//		log.Println("Ошибка при сохранении рецепта в Redis:", err)
//	}
//}
//
//func handleRecipeConfirmationState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	if update.CallbackQuery != nil {
//		switch update.CallbackQuery.Data {
//		case "recipe_confirm":
//			reply := "Рецепт успешно создан!"
//			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
//			bot.Send(msg)
//
//			recipe.RecipeState = StateRecipeIdle
//			err := redisClient.Del(context.Background(), fmt.Sprintf("recipe:%d", recipe.UserID)).Err()
//			if err != nil {
//				log.Println("Ошибка при удалении данных рецепта из Redis:", err)
//			}
//
//		case "recipe_restart":
//			reply := "Процесс создания рецепта будет начат заново."
//			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
//			bot.Send(msg)
//
//			recipe.RecipeState = StateRecipeIdle
//			err := redisClient.Del(context.Background(), fmt.Sprintf("recipe:%d", recipe.UserID)).Err()
//			if err != nil {
//				log.Println("Ошибка при удалении данных рецепта из Redis:", err)
//			}
//		}
//	}
//}
