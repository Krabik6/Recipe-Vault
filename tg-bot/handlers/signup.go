package handlers

//
//import (
//	"fmt"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/redis/go-redis/v9"
//	"log"
//	"net/http"
//	"strings"
//)
//
//func handleIdleState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	if update.Message.Text == "/signUp" {
//		reply := "Добро пожаловать в процесс регистрации!\nПожалуйста, введите свой логин:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//		user.RegStatus = StateWaitingForLogin
//		user.UserID = update.Message.Chat.ID // Сохраняем UserID пользователя
//		err := saveUserToRedis(redisClient, user)
//		if err != nil {
//			log.Println("Ошибка при сохранении пользователя в Redis:", err)
//		}
//	} else if update.Message.Text == "/createRecipe" {
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
//func handleLoginState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	login := strings.TrimSpace(update.Message.Text)
//	if login == "" {
//		reply := "Логин не может быть пустым. Пожалуйста, введите логин:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//		return
//	}
//
//	user.Login = login
//
//	reply := "Отлично! Теперь введите свой никнейм:"
//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//	bot.Send(msg)
//
//	user.RegStatus = StateWaitingForNickname
//	err := saveUserToRedis(redisClient, user)
//	if err != nil {
//		log.Println("Ошибка при сохранении пользователя в Redis:", err)
//	}
//}
//
//func handleNicknameState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	nickname := strings.TrimSpace(update.Message.Text)
//	if nickname == "" {
//		reply := "Никнейм не может быть пустым. Пожалуйста, введите никнейм:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//		return
//	}
//
//	user.Nickname = nickname
//
//	reply := "Отлично! Теперь введите свой пароль:"
//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//	bot.Send(msg)
//
//	user.RegStatus = StateWaitingForPassword
//	err := saveUserToRedis(redisClient, user)
//	if err != nil {
//		log.Println("Ошибка при сохранении пользователя в Redis:", err)
//	}
//}
//
//func handlePasswordState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	password := strings.TrimSpace(update.Message.Text)
//	if password == "" {
//		reply := "Пароль не может быть пустым. Пожалуйста, введите пароль:"
//		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//		bot.Send(msg)
//		return
//	}
//
//	user.Password = password
//
//	reply := fmt.Sprintf("Вы ввели следующие данные:\nЛогин: %s\nНикнейм: %s\nПароль: %s\n\nПожалуйста, подтвердите, что все верно.",
//		user.Login, user.Nickname, user.Password)
//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
//	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("Да", "confirm"),
//			tgbotapi.NewInlineKeyboardButtonData("Нет", "restart"),
//		),
//	)
//	bot.Send(msg)
//
//	user.RegStatus = StateConfirmation
//	err := saveUserToRedis(redisClient, user)
//	if err != nil {
//		log.Println("Ошибка при сохранении пользователя в Redis:", err)
//	}
//}
//
//func handleConfirmationState(bot *tgbotapi.BotAPI, redisClient *redis.Client, user *User, recipe *Recipe, update tgbotapi.Update) {
//	if update.CallbackQuery != nil {
//		switch update.CallbackQuery.Data {
//		case "confirm":
//			client := &http.Client{}
//			SignUp(bot, update, client, *user)
//			reply := "Регистрация успешно завершена!"
//			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
//			bot.Send(msg)
//
//			user.RegStatus = StateIdle
//			err := redisClient.Del(context.Background(), fmt.Sprintf("user:%d", user.UserID)).Err()
//			if err != nil {
//				log.Println("Ошибка при удалении данных пользователя из Redis:", err)
//			}
//
//		case "restart":
//			reply := "Процесс регистрации будет начат заново."
//			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
//			bot.Send(msg)
//
//			user.RegStatus = StateIdle
//			err := redisClient.Del(context.Background(), fmt.Sprintf("user:%d", user.UserID)).Err()
//			if err != nil {
//				log.Println("Ошибка при удалении данных пользователя из Redis:", err)
//			}
//		case "cancel":
//			handleCancelState(bot, redisClient, user, recipe, update)
//		}
//	}
//}
