package interfaces

import "context"

type JwtManager interface {
	SetUserJWTToken(ctx context.Context, userID int64, token string) error
	GetUserJWTToken(ctx context.Context, userID int64) (string, error)
	DeleteUserJWTToken(ctx context.Context, userID int64) error
	CheckLoggedIn(ctx context.Context, userID int64) (bool, error)
}
