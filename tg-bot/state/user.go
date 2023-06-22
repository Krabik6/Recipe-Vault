package state

// struct user, that goona be get from Redis by user id and store in Redis by user id
type User struct {
	GlobalState         GlobalState
	RegistrationState   RegistrationState
	RecipeCreationState RecipeCreationState
}
