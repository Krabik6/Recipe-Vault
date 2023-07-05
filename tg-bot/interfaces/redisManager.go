package interfaces

import (
	"context"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
)

//type RedisManager struct {
//	UserStateManager UserStateManager
//}

//type LocalState int

// UserStateManager defines the methods for managing user states in Redis.
type UserStateManager interface {
	SetUserState(ctx context.Context, userID int64, state model.State) error
	GetUserState(ctx context.Context, userID int64) (model.State, error)
	DeleteUserState(ctx context.Context, userID int64) error
}

//// UserDataManager defines the methods for managing user-specific data in Redis.
//type UserDataManager interface {
//	SetUserData(ctx context.Context, userID int64, data UserData) error
//	GetUserData(ctx context.Context, userID int64) (UserData, error)
//	DeleteUserData(ctx context.Context, userID int64) error
//}
