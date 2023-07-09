package registration

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/api"
	"github.com/Krabik6/meal-schedule/tg-bot/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
)

type Registration interface {
	SetUserRegistrationState(ctx context.Context, userID int64) error
	DeleteUserRegistrationState(ctx context.Context, userID int64) error
	GetUserRegistrationState(ctx context.Context, userID int64) (registrationState, error)
	SetUserRegistrationData(ctx context.Context, userID int64) error
	DeleteUserRegistrationData(ctx context.Context, userID int64) error
	GetUserRegistrationData(ctx context.Context, userID int64) (string, string, string, error)
	BuildUserData(name string, email string, password string)
	HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error
	HandleState(ctx context.Context, userID int64, update tgbotapi.Update) error
	HandleMessage(ctx context.Context, userID int64, message string, update tgbotapi.Update) error
}

type RegistrationService struct {
	Registration
}

func NewRegistrationService(
	bot *tgbotapi.BotAPI,
	client *redis.Client,
	userManager interfaces.UserStateManager,
	botMenu interfaces.BotMenu,
	api *api.Api,
) *RegistrationService {
	return &RegistrationService{
		Registration: NewRegistrationStateHandler(
			client,
			userManager,
			bot,
			botMenu,
			api,
		),
	}
}
