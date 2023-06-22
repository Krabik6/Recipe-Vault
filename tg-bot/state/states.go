package state

// variable global state to store the state for user (registration, recipe creation, etc.) in numeric format, using next idea
// 0 - no state
// 1 - registration
// 2 - recipe creation

type GlobalState int

const (
	NoState GlobalState = iota
	Registration
	RecipeCreation
)

type RegistrationState int

const (
	NoRegistrationState RegistrationState = iota
	RegistrationName
	RegistrationEmail
	RegistrationPassword
	RegistrationConfirmPassword
	RegistrationComplete
)

type RecipeCreationState int

const (
	NoRecipeCreationState RecipeCreationState = iota
	RecipeCreationName
	RecipeCreationDescription
	RecipeCreationIngredients
	RecipeCreationComplete
)
