package service

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"strings"
)

func SignUp(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sign up")
	bot.Send(msg)

	args := update.Message.CommandArguments()
	if args == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your username and password in the format /signup <username> <password> <name>."))
		return
	}
	argList := strings.Split(args, " ")
	if len(argList) != 3 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your username and password and name in the format /signup <username> <password> <name>."))
		return
	}
	username := argList[0]
	password := argList[1]
	name := argList[2]

	signUpCredentials := model.SignUpCredentials{
		Username: username,
		Password: password,
		Name:     name,
	}

	requestBody, err := json.Marshal(signUpCredentials)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing up ( \"+err.Error()+\" ). Please try again."))
		return
	}

	resp, err := client.Post("http://localhost:8000/auth/sign-up", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing up ( "+err.Error()+" ). Please try again. POST"))
		return
	}
	defer resp.Body.Close()

	var authResponse model.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while signing up ( "+err.Error()+" ). Please try again."))
		return
	}

	fmt.Println(resp.StatusCode, authResponse.Token)

	message := "You have successfully signed up. Now signIn" + authResponse.Token
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
}
