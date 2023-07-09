package login

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type Login interface {
	BuildUserData(userName string, password string)
	HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error
	HandleState(ctx context.Context, userID int64, update tgbotapi.Update) error
	HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error
	SetUserLoginState(ctx context.Context, userID int64) error
	DeleteUserLoginState(ctx context.Context, userID int64) error
	GetUserLoginState(ctx context.Context, userID int64) (loginState, error)
	SetUserLoginData(ctx context.Context, userID int64) error
	DeleteUserLoginData(ctx context.Context, userID int64) error
	GetUserLoginData(ctx context.Context, userID int64) (string, string, error)
}

type LoginService struct {
	Login
}

func NewLoginService(
	client *redis.Client,
	bot *tgbotapi.BotAPI,
	jwtManager interfaces.JwtManager,
	userStateManager interfaces.UserStateManager,
	botMenu interfaces.BotMenu,
	api *api.Api,
) *LoginService {
	return &LoginService{
		Login: NewLoginStateHandler(
			client,
			bot,
			jwtManager,
			userStateManager,
			botMenu,
			api),
	}
}
