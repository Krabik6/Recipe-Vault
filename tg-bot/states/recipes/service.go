package recipes

// struct for recipes methods

/*
interface for
Methods on (*CreateRecipeStateHandler):
SetUserState(ctx context.Context, userID int64) error
DeleteUserState(ctx context.Context, userID int64) error
GetUserState(ctx context.Context, userID int64) (createRecipeState, error)
SetUserData(ctx context.Context, userID int64) error
DeleteUserData(ctx context.Context, userID int64) error
GetUserData(ctx context.Context, userID int64) (title string, description string, isPublic bool, cost float64, timeToPrepare int64, healthy int, err error)
HandleMessage(ctx context.Context, userID int64, update tgbotapi.Update) error
HandleCallbackQuery(ctx context.Context, userID int64, update tgbotapi.Update) error
handleState(ctx context.Context, userID int64) error
handleNoCreateRecipeState(userID int64) error
handleCreateRecipeTitle(userID int64) error
handleCreateRecipeDescription(userID int64) error
handleCreateRecipeIsPublic(userID int64) error
handleCreateRecipeCost(userID int64) error
handleCreateRecipeTimeToPrepare(userID int64) error
handleCreateRecipeHealthy(userID int64) error
handleCreateRecipeConfirmYes(ctx context.Context, userID int64) error
handleCreateRecipeConfirmNo(ctx context.Context, userID int64) error
handleCancel(ctx context.Context, userID int64) error

*/

type CreateRecipeStateHandler interface {
}

type Service struct {
}
