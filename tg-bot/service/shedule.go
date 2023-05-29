package service

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/tg-bot/cache"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"github.com/Krabik6/meal-schedule/tg-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func CreateMeal(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository, baseUrl string) {
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
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your meal in the format /CreateMeal <name> <at_time(yyyy-mm-dd hh:mm:ss)> <[recipe ids]>."))
		return
	}
	fmt.Println(args, "args")
	argList := strings.Split(args, "  ")
	if len(argList) < 3 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Not enough args /CreateMeal <name> <at_time(yyyy-mm-dd hh:mm:ss)> <[1,2,3,4,5]>."))
		return
	}

	name := argList[0]
	atTime := argList[1]
	recipeIdsStr := argList[2]
	recipeIds := strings.Split(recipeIdsStr, ",")
	fmt.Println(recipeIds, "recipeIds")

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

	req, err := http.NewRequest("POST", "http://localhost:8000/api/schedule/meal", strings.NewReader(string(requestBody)))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	defer req.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error, with status not found. Code: %d", resp.StatusCode)))
		return
	} else if resp.StatusCode != http.StatusOK {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error while creating meal. Status code: %d", resp.StatusCode)))
		return
	}

	var scheduleResponse model.ScheduleResponse
	err = json.NewDecoder(resp.Body).Decode(&scheduleResponse)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating meal ( "+err.Error()+" ). Please try again."))
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("You have successfully added a meal with id %d", scheduleResponse.Id))
	bot.Send(msg)

}

func GetScheduleByPeriod(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository, baseURl string) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	args := update.Message.CommandArguments()
	if args == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your strDate in the format /getScheduleByDate <strDate>  <strPeriod>."))
		return
	}

	argList := strings.Split(args, " ")
	if len(argList) < 2 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your strDate in the format /getScheduleByDate <strDate>  <strPeriod> in format."))
		return
	}

	strDate := argList[0]
	strPeriod := argList[1]

	log.Println("strDate: ", strDate)
	log.Println("strPeriod: ", strPeriod)

	stringReq := baseURl + "/api/schedule/?date=" + strDate + "&period=" + strPeriod

	req, err := http.NewRequest("GET", stringReq, nil)
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

	if resp.StatusCode == http.StatusNotFound {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error, with status not found. Code: %d", resp.StatusCode)))
		return
	} else if resp.StatusCode != http.StatusOK {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error while creating meal. Status code: %d", resp.StatusCode)))
		return
	}

	var meal []models.ScheduleByDateOutput
	err = json.NewDecoder(resp.Body).Decode(&meal)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting schedule ( "+err.Error()+" ). Please try again."))
		return
	}

	date, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting schedule ( "+err.Error()+" ). Please try again."))
		return
	}

	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting schedule ( "+err.Error()+" ). Please try again."))
		return
	}
	date2 := date.AddDate(0, 0, period)

	layout := "2006-01-02"

	strMsg := fmt.Sprintf("Schedule from %s to %s:\n\n%s", date.Format(layout), date2.Format(layout), utils.ScheduleMealsOutputToString(meal))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strMsg)
	bot.Send(msg)

}
