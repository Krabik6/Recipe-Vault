package interfaces

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
)

type Notate interface {
	HandleMessage(ctx context.Context, userID int64, command string, state model.State) (model.State, error)
	UnknownCommand(ctx context.Context, userID int64) error
	Start(ctx context.Context, userID int64) error
	Help(ctx context.Context, userID int64) error
	HandleCallback(ctx context.Context, userID int64, callbackData string) error
	LogOut(ctx context.Context, userID int64) error
}
