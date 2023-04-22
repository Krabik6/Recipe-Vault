package service

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/tg-bot/cache"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"github.com/Krabik6/meal-schedule/tg-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func FillSchedule(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	args := update.Message.CommandArguments()
	if args == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your meal in the format /addMeal <name> <at_time> <[recipe ids]>."))
		return
	}

	argList := strings.Split(args, " ")
	if len(argList) < 3 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your meal in the format /addMeal <name> <at_time> <[1,2,3,4,5]>."))
		return
	}

	name := argList[0]
	atTime := argList[1]
	recipeIdsStr := argList[2]
	recipeIds := strings.Split(recipeIdsStr, ",")

	//array of strings to array of int
	//var recipeIdsInt []int
	//for _, v := range recipeIds {
	//	recipeIdsInt = append(recipeIdsInt, strconv.Atoi(v))
	//}
	recipes, err := utils.StringArrayToIntArray(recipeIds)
	meal := models.Meal{
		Name:    name,
		AtTime:  atTime,
		Recipes: recipes,
	}
	fmt.Println(meal, "meal")

	requestBody, err := json.Marshal(meal)
	if err != nil {
		fmt.Println("Error occurred:")
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while filling schedule ( "+err.Error()+" ). Please try again."))
		return
	}

	//req, err := http.NewRequest("POST", "http://localhost:8000/schedule/", strings.NewReader(string(requestBody)))
	//if err != nil {
	//	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while filling schedule ( "+err.Error()+" ). Please try again."))
	//	return
	//}
	//req.Header.Set("Authorization", "Bearer "+token)
	//
	resp, err := client.Post("http://localhost:8000/schedule/", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		fmt.Println("Error occurred:")
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while filling schedule ( "+err.Error()+" ). Please try again."))
		return
	}
	resp.Header.Set("Authorization", "Bearer "+token)

	defer resp.Body.Close()

	var scheduleResponse model.ScheduleResponse
	err = json.NewDecoder(resp.Body).Decode(&scheduleResponse.Id)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while filling schedule ( "+err.Error()+" ). Please try again."))
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("You have successfully added a meal with id %d", scheduleResponse.Id))
	bot.Send(msg)

}

func GetScheduleByDate(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	args := update.Message.CommandArguments()
	if args == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your date in the format /getScheduleByDate <date>."))
		return
	}

	argList := strings.Split(args, " ")
	if len(argList) < 1 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your date in the format /getScheduleByDate <date>."))
		return
	}

	date := argList[0]

	req, err := http.NewRequest("GET", "http://localhost:8000/schedule/?date="+date, nil)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting schedule ( "+err.Error()+" ). Please try again."))
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting schedule ( "+err.Error()+" ). Please try again."))
		return
	}
	defer resp.Body.Close()
	/*
	   Id
	   Date
	   BreakfastI
	   LunchId
	   DinnerId
	*/
	var meal []models.Meal
	err = json.NewDecoder(resp.Body).Decode(&meal)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting schedule ( "+err.Error()+" ). Please try again."))
		return
	}

	//msg := tgbotapi.NewMessage(update.Message.Chat.ID,)
	// msg with schedule response

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, utils.ArrayToString(meal))
	bot.Send(msg)

}
