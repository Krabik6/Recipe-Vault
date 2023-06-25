package api

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strings"
)

func SignUp(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, user model.SignUpCredentials) error {
	signUpCredentials := model.SignUpCredentials{
		Username: user.Username,
		Password: user.Password,
		Name:     user.Name,
	}

	requestBody, err := json.Marshal(signUpCredentials)
	if err != nil {
		return err
	}

	resp, err := client.Post("http://localhost:8000/auth/sign-up", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("status code is %d. \n response: %s", resp.StatusCode, body)
	}

	return nil
}
