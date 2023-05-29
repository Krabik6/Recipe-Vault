package service

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/cache"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

const BaseURL = "http://localhost:8000"

func SignIn(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository) {
	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sign in")
	//bot.Send(msg)

	args := update.Message.CommandArguments()
	if args == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your username and password in the format /signin <username> <password>."))
		return
	}
	argList := strings.Split(args, " ")
	if len(argList) != 2 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your username and password in the format /signin <username> <password>."))
		return
	}
	username := argList[0]
	password := argList[1]

	signInCredentials := model.SignInCredentials{
		Username: username,
		Password: password,
	}

	requestBody, err := json.Marshal(signInCredentials)
	if err != nil {
		debug.PrintStack()
		_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing in ( \"+err.Error()+\" ). Please try again."))
		if err != nil {
			return
		}

		return
	}

	resp, err := client.Post("http://localhost:8000/auth/sign-in", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing in ( "+err.Error()+" ). Please try again. POST"))
		debug.PrintStack()
		return
	}
	defer resp.Body.Close()

	var authResponse model.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing in ( "+err.Error()+" ). Please try again."))
		return
	}

	err = CRepo.SetKey(strconv.FormatInt(update.Message.Chat.ID, 10), authResponse.Token, 0)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing in ( "+err.Error()+" ). Please try again."))
		return
	}
	fmt.Println(CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10)))

	message := "You have successfully signed in. Your token is: " + authResponse.Token
	_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	if err != nil {
		return
	}
}
