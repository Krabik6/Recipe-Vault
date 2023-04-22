package main

import (
	"flag"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/cache"
	"github.com/Krabik6/meal-schedule/tg-bot/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

const (

	// 10 lines of random text for the bot
	replyStart                 = `Hello, I am a bot that can help you with your meals.`
	replyUnknown               = "I don't know what you mean"
	replyError                 = "Something went wrong"
	replySuccess               = "Success"
	replyEmpty                 = "Empty"
	replyCancel                = "Cancel"
	replyMealAdded             = "Meal added"
	replyMealDeleted           = "Meal deleted"
	replyMealUpdated           = "Meal updated"
	replyGetMeals              = "Meals"
	replyGetMealsByDate        = "Meals by date"
	replyGetMealsByDateAndTime = "Meals by date and time"
	replyGetRecipes            = "Recipes"
	replyGetIngredients        = "Ingredients"
	replyHelp                  = "Help menu for the bot"
)
const BaseURL = "http://localhost:8000"

//func that reply to /help, its a list of commands (signIn
//signUp
//createRecipe
//getAllRecipes
//getRecipeById
//updateRecipe
//getFilteredRecipes
//getFilteredUserRecipes
//deleteRecipe)

func replyHelpMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyHelp)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/signIn"),
			tgbotapi.NewKeyboardButton("/signUp"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/createRecipe"),
			tgbotapi.NewKeyboardButton("/getAllRecipes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/getRecipeById"),
			tgbotapi.NewKeyboardButton("/updateRecipe"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/getFilteredRecipes"),
			tgbotapi.NewKeyboardButton("/getFilteredUserRecipes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/deleteRecipe"),
		),
	)
	bot.Send(msg)
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	CRepo := cache.NewCacheRepository(redisClient)
	//err := CRepo.CRD.SetKey("hello", "world", 0)
	//if err != nil {
	//	panic(err)
	//}
	//val, err := CRepo.CRD.GetKey("hello")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(val)

	bot, err := tgbotapi.NewBotAPI(mustToken())
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	bot.Debug = true
	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates := bot.GetUpdatesChan(ucfg)

	for {
		select {
		case update := <-updates:
			//
			//UserName := update.Message.From.UserName
			//
			//ChatId := update.Message.Chat.ID
			//
			//Text := update.Message.Text

			if update.Message == nil {
				// Ignore non-message updates
				continue
			}

			if update.Message.Command() == "signIn" {
				service.SignIn(bot, update, client, CRepo)
			}

			if update.Message.Command() == "signUp" {
				service.SignUp(bot, update, client)
			}
			// addMeal <name> <date> <time> <recipe> <ingredients>
			if update.Message.Command() == "addMeal" {
				service.FillSchedule(bot, update, client, CRepo)
			}

			if update.Message.Command() == "getMeals" {
				service.GetScheduleByDate(bot, update, client, CRepo)
			}

			if update.Message.Text == "/start" {
				replyToMessage(bot, update, replyStart)
			}

			if update.Message.Text == "/help" {
				replyHelpMenu(bot, update)
			}

			// createRecipe
			if update.Message.Command() == "createRecipe" {
				service.CreateRecipe(bot, update, client, CRepo)
			}
			//GetAllRecipes
			if update.Message.Command() == "getAllRecipes" {
				service.GetAllRecipes(bot, update, client, CRepo)
			}

			//GetRecipeById
			if update.Message.Command() == "getRecipeById" {
				service.GetRecipeById(bot, update, client, CRepo)
			}

			// updateRecipe
			if update.Message.Command() == "updateRecipe" {
				service.UpdateRecipe(bot, update, client, CRepo, BaseURL)
			}

			// GetFilteredRecipes
			if update.Message.Command() == "getFilteredRecipes" {
				service.GetFilteredRecipes(bot, update, client, CRepo, BaseURL)
			}

			//GetFilteredUserRecipes
			if update.Message.Command() == "getFilteredUserRecipes" {
				service.GetFilteredUserRecipes(bot, update, client, CRepo, BaseURL)
			}

			//deleteRecipe
			if update.Message.Command() == "deleteRecipe" {
				service.DeleteRecipe(bot, update, client, CRepo, BaseURL)
			}

			//replyToMessage(bot, update, replyStart)

			//if Text == "Hello world" {
			//	text := "```package main \n import fmt \n func main() { \n fmt.Println(\"Hello world\") \n }```"
			//	msg := tgbotapi.NewMessage(ChatId, text)
			//	msg.ParseMode = "markdown"
			//	bot.Send(msg)
			//	continue
			//}

		}
	}

}

func mustToken() string {
	token := flag.String("token", "", "Telegram bot token")

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is required")
	}
	return *token
}

func replyToMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

/*
	schedule := api.Group("/schedule")
		{
			schedule.POST("/", h.fillSchedule)
			schedule.POST("/meal", h.createMeal)
			schedule.GET("/all", h.getAllSchedule)
			schedule.GET("/", h.getScheduleByDate)
			schedule.PUT("/", h.updateSchedule)
			schedule.DELETE("/", h.deleteSchedule)
		}
	FillSchedule(userId int, meal models.Meal) (int, error)

*/
