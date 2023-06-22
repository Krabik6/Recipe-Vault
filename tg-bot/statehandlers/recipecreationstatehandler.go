package statehandlers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"runtime/debug"
)

type RecipeCreationStateHandler struct {
	State  recipeCreationState
	Client *redis.Client
	Recipe Recipe
}

type recipeCreationState int

const (
	NoRecipeCreationState recipeCreationState = iota
	RecipeCreationName
	RecipeCreationIngredients
	RecipeCreationInstructions
	RecipeCreationComplete
)

type Recipe struct {
	Name         string
	Ingredients  []string
	Instructions string
}

func (h *RecipeCreationStateHandler) HandleState(ctx context.Context, userID int64) error {
	switch h.State {
	case NoRecipeCreationState:
		return h.handleNoRecipeCreationState(ctx, userID)
	case RecipeCreationName:
		return h.handleRecipeCreationName(ctx, userID)
	case RecipeCreationIngredients:
		return h.handleRecipeCreationIngredients(ctx, userID)
	case RecipeCreationInstructions:
		return h.handleRecipeCreationInstructions(ctx, userID)
	case RecipeCreationComplete:
		return h.handleRecipeCreationComplete(ctx, userID)
	default:
		return fmt.Errorf("unknown recipe creation state")
	}
}

func (h *RecipeCreationStateHandler) HandleMessage(ctx context.Context, userID int64, message string) error {
	switch h.State {
	case RecipeCreationName:
		h.Recipe.Name = message
	case RecipeCreationIngredients:
		h.Recipe.Ingredients = append(h.Recipe.Ingredients, message)
	case RecipeCreationInstructions:
		h.Recipe.Instructions = message
	default:
		return fmt.Errorf("unexpected message received in current state")
	}

	// Обработка состояния
	err := h.HandleState(ctx, userID)
	if err != nil {
		return err
	}

	// Вывод текущего состояния
	fmt.Println(h.State)

	return nil
}

func (h *RecipeCreationStateHandler) handleNoRecipeCreationState(ctx context.Context, userID int64) error {
	fmt.Println("Welcome! Please provide the name of the recipe.")
	// Меняем состояние создания рецепта на ввод имени
	h.State = RecipeCreationName

	// Сохраняем состояние в Redis
	return h.saveStateToRedis(ctx, userID)
}

func (h *RecipeCreationStateHandler) handleRecipeCreationName(ctx context.Context, userID int64) error {
	fmt.Println("Please provide the list of ingredients for the recipe.")
	// Меняем состояние создания рецепта на ввод ингредиентов
	h.State = RecipeCreationIngredients

	// Сохраняем состояние в Redis
	return h.saveStateToRedis(ctx, userID)
}

func (h *RecipeCreationStateHandler) handleRecipeCreationIngredients(ctx context.Context, userID int64) error {
	fmt.Println("Please provide the instructions for the recipe.")
	// Меняем состояние создания рецепта на ввод инструкций
	h.State = RecipeCreationInstructions

	// Сохраняем состояние в Redis
	return h.saveStateToRedis(ctx, userID)
}

func (h *RecipeCreationStateHandler) handleRecipeCreationInstructions(ctx context.Context, userID int64) error {
	fmt.Println("Recipe creation complete!")
	// Меняем состояние создания рецепта на завершено
	h.State = RecipeCreationComplete

	// Сохраняем состояние в Redis
	return h.saveStateToRedis(ctx, userID)
}

func (h *RecipeCreationStateHandler) handleRecipeCreationComplete(ctx context.Context, userID int64) error {
	// Дополнительная логика после завершения создания рецепта
	fmt.Println("Recipe created successfully!")

	// Удаляем состояние из Redis, так как создание рецепта завершено
	return h.deleteStateFromRedis(ctx, userID)
}

func (h *RecipeCreationStateHandler) saveStateToRedis(ctx context.Context, userID int64) error {
	err := h.Client.Set(ctx, fmt.Sprintf("state:%d", userID), h.State, 0).Err()
	if err != nil {
		return err
	}

	// Сохраняем данные рецепта в Redis
	err = h.Client.HSet(ctx, fmt.Sprintf("recipe:%d", userID), "name", h.Recipe.Name, "ingredients", h.Recipe.Ingredients, "instructions", h.Recipe.Instructions).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *RecipeCreationStateHandler) deleteStateFromRedis(ctx context.Context, userID int64) error {
	err := h.Client.Del(ctx, fmt.Sprintf("state:%d", userID)).Err()
	if err != nil {
		return err
	}

	// Удаляем данные рецепта из Redis
	err = h.Client.Del(ctx, fmt.Sprintf("recipe:%d", userID)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *RecipeCreationStateHandler) HandleCommand(ctx context.Context, userID int64, input string) error {
	switch h.State {
	case RecipeCreationName:
		h.Recipe.Name = input
	case RecipeCreationIngredients:
		h.Recipe.Ingredients = append(h.Recipe.Ingredients, input)
	case RecipeCreationInstructions:
		h.Recipe.Instructions = input

	default:
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return fmt.Errorf("unknown registration state")
	}

	// Сохранение состояния в Redis
	err := h.saveStateToRedis(ctx, userID)
	if err != nil {
		stackTrace := debug.Stack()
		fmt.Println("Stack trace:", string(stackTrace))
		return err
	}

	// Вывод текущего состояния
	fmt.Println(h.State)

	return nil
}
