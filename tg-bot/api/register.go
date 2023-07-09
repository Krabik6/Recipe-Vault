package api

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	"io"
	"net/http"
	"strings"
)

func (a *Api) SignUp(user model.SignUpCredentials) error {
	signUpCredentials := model.SignUpCredentials{
		Username: user.Username,
		Password: user.Password,
		Name:     user.Name,
	}

	requestBody, err := json.Marshal(signUpCredentials)
	if err != nil {
		return err
	}
	//todo replace hardcoded url to env variable
	resp, err := a.Client.Post("http://localhost:8000/auth/sign-up", "application/json", strings.NewReader(string(requestBody)))
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
