package service

import (
	json "encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/tg-bot/cache"
	"github.com/Krabik6/meal-schedule/tg-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
)

func CreateRecipe(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	args := update.Message.CommandArguments()
	if args == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your recipe in the format /createRecipe <title> <description> <is_public> <cost> <time_to_prepare> <healthy>"))
		return
	}

	argList := strings.Split(args, " ")
	if len(argList) < 5 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your recipe in the format /createRecipe <title> <description> <is_public> <cost> <time_to_prepare> <healthy>"))
		return
	}

	title := argList[0]
	description := argList[1]
	isPublic := argList[2]
	cost := argList[3]
	timeToPrepare := argList[4]
	healthy := argList[5]

	//convert isPublic to bool
	isPublicBool, err := strconv.ParseBool(isPublic)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	// convert cost to float
	costFloat, err := strconv.ParseFloat(cost, 64)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	// convert timeToPrepare to int64
	timeToPrepareInt, err := strconv.ParseInt(timeToPrepare, 10, 64)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	// convert healthy to f
	healthyInt, err := strconv.Atoi(healthy)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	recipe := models.Recipe{
		Title:         title,
		Description:   description,
		IsPublic:      isPublicBool,
		Cost:          costFloat,
		TimeToPrepare: timeToPrepareInt,
		Healthy:       healthyInt,
	}
	fmt.Println(recipe, "recipe")

	requestBody, err := json.Marshal(recipe)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/api/recipes/", strings.NewReader(string(requestBody)))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	var recipeResp model.CreateRecipeResponse
	err = json.NewDecoder(resp.Body).Decode(&recipeResp)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while creating recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Recipe created successfully. Id: "+strconv.FormatInt(recipeResp.Id, 10)+""))

}

/*
api := router.Group("/api", h.userIdentity)
	{
		recipes := api.Group("/recipes")
		{
			recipes.POST("/", h.createRecipe)
			recipes.GET("/", h.getAllRecipes)
			recipes.GET("/:id", h.getRecipeById)
			recipes.PUT("/:id", h.updateRecipe)
			recipes.DELETE("/:id", h.deleteRecipe)
			recipes.GET("/public", h.getPublicRecipes)
		}
*/

func GetAllRecipes(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8000/api/recipes/", nil)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Recipe not found."))
		return
	} else if resp.StatusCode != http.StatusOK {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error while getting recipe. %d", resp.StatusCode)))
		return
	}

	var recipes []models.Recipe
	err = json.NewDecoder(resp.Body).Decode(&recipes)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	var messageText string
	if len(recipes) == 0 {
		messageText = "No recipes found."
	} else {
		for _, recipe := range recipes {
			messageText += fmt.Sprintf("Id: %d\nTitle: %s\nDescription: %s\nCost: %.2f\nTime to Prepare: %d minutes\nHealthy: %d\n\n", recipe.Id, recipe.Title, recipe.Description, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	bot.Send(msg)
}

func GetRecipeById(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	args := update.Message.CommandArguments()
	if args == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide the recipe ID in the format /getRecipeById <id>."))
		return
	}

	recipeID, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid recipe ID. Please provide a valid integer ID."))
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8000/api/recipes/"+strconv.FormatInt(recipeID, 10), nil)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Recipe not found."))
		return
	} else if resp.StatusCode != http.StatusOK {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe."))
		return
	}

	var recipe models.Recipe
	err = json.NewDecoder(resp.Body).Decode(&recipe)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	messageText := fmt.Sprintf("Id: %d\n Title:%s\nDescription: %s\nCost: %.2f\nTime to Prepare: %d minutes\nHealthy: %d", recipe.Id, recipe.Title, recipe.Description, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	bot.Send(msg)
}

// UpdateRecipe func updateRecipe that have omitempty fields
/*

type RecipesFilter struct {
	CostMoreThan          *float64 `json:"costMoreThan,omitempty"`
	CostLessThan          *float64 `json:"costLessThan,omitempty"`
	TimeToPrepareMoreThan *int     `json:"timeToPrepareMoreThan,omitempty"`
	TimeToPrepareLessThan *int     `json:"timeToPrepareLessThan,omitempty"`
	HealthyMoreThan       *int     `json:"healthyMoreThan,omitempty"`
	HealthyLessThan       *int     `json:"healthyLessThan,omitempty"`
}

type UpdateRecipeInput struct {
	Id            *int     `json:"id,omitempty" db:"id"`
	Title         *string  `json:"title"`
	Description   *string  `json:"description"`
	IsPublic      *bool    `json:"isPublic" db:"public"`
	Cost          *float64 `json:"cost,omitempty"`
	TimeToPrepare *int     `json:"timeToPrepare,omitempty"`
	Healthy       *int     `json:"healthy,omitempty"`
}
*/
func UpdateRecipe(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository, baseUrl string) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	args := update.Message.CommandArguments()
	if args == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your recipe update in the format  /updateRecipe <id> <Title> <Description> <IsPublic> <Cost> <TimeToPrepare> <Healthy>. If some fields are not needed, please provide \"-\"."))
		return
	}

	argList := strings.Split(args, " ")
	if len(argList) < 7 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your recipe update in the format  /updateRecipe <id> <Title> <Description> <IsPublic> <Cost> <TimeToPrepare> <Healthy>. If some fields are not needed, please provide \"-\"."))
		return
	}

	recipeID, err := strconv.ParseInt(argList[0], 10, 64)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid recipe ID. Please provide a valid integer ID."))
		return
	}

	id := int(recipeID)
	title := argList[1]
	description := argList[2]
	isPublic := argList[3]
	cost := argList[4]
	timeToPrepare := argList[5]
	healthy := argList[6]

	var recipeUpdate models.UpdateRecipeInput
	recipeUpdate.Id = &id
	if title != "-" {
		recipeUpdate.Title = &title
	}

	if description != "-" {
		recipeUpdate.Description = &description
	}

	if isPublic != "-" {
		isPublicBool, err := strconv.ParseBool(isPublic)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid isPublic value. Please provide a valid boolean value."))
			return
		}
		recipeUpdate.IsPublic = &isPublicBool
	}

	if cost != "-" {
		costFloat, err := strconv.ParseFloat(cost, 64)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid cost value. Please provide a valid float value."))
			return
		}
		recipeUpdate.Cost = &costFloat
	}

	if timeToPrepare != "-" {
		timeToPrepareInt, err := strconv.Atoi(timeToPrepare)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid timeToPrepare value. Please provide a valid integer value."))
			return
		}
		timeToPrepareInt = timeToPrepareInt * 60
		recipeUpdate.TimeToPrepare = &timeToPrepareInt
	}

	//healthy enum 1,2,3
	if healthy != "-" {
		healthyInt, err := strconv.Atoi(healthy)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value."))
			return
		}
		if healthyInt < 1 || healthyInt > 3 {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value (Healthy only 1, 2, 3)."))
			return
		}
		recipeUpdate.Healthy = &healthyInt
	}

	//if all fields are "-", return error
	if recipeUpdate.Title == nil && recipeUpdate.Description == nil && recipeUpdate.IsPublic == nil && recipeUpdate.Cost == nil && recipeUpdate.TimeToPrepare == nil && recipeUpdate.Healthy == nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your recipe update in the format  /updateRecipe <id> <Title> <Description> <IsPublic> <Cost> <TimeToPrepare> <Healthy>. If some fields are not needed, please provide \"-\"."))
		return
	}

	//params list, only for update, theyre gonna put in request after "?" in url
	stringRecipeId := strconv.FormatInt(recipeID, 10)
	stringRequest := baseUrl + "/api/recipes/" + stringRecipeId
	log.Println(stringRequest)

	//i need to put in request body, so i need to marshal it
	jsonRecipe, err := json.Marshal(recipeUpdate)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req, err := http.NewRequest("PUT", stringRequest, strings.NewReader(string(jsonRecipe)))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error, with status not found. Code: %d", resp.StatusCode)))
		return
	} else if resp.StatusCode != http.StatusOK {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error while updating recipe. Status code: %d", resp.StatusCode)))
		return
	}

	var response model.UpdateRecipeResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while updating recipe ( "+err.Error()+" ). Please try again."))
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Recipe %d updated successfully. ", recipeID)))
}

func GetFilteredRecipes(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository, baseUrl string) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	//get args
	args := strings.Split(update.Message.Text, " ")
	fmt.Println(len(args))
	if len(args) < 7 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your filter in the format  /getFilteredRecipes <costMoreThat> <costLessThat> <timeToPrepareMoreThan> <timeToPrepareLessThan> <healthyMoreThan> <healthyLessThan> . If some fields are not needed, please provide \"-\"."))
		return
	}

	costMoreThat := args[1]
	costLessThat := args[2]
	timeToPrepareMoreThan := args[3]
	timeToPrepareLessThan := args[4]
	healthyMoreThan := args[5]
	healthyLessThan := args[6]
	params := url.Values{}

	//check if args are valid and add them to params
	if costMoreThat != "-" {
		params.Add("costMoreThat", costMoreThat)
	}

	if costLessThat != "-" {
		params.Add("costLessThat", costLessThat)
	}

	if timeToPrepareMoreThan != "-" {
		params.Add("timeToPrepareMoreThan", timeToPrepareMoreThan)
	}

	if timeToPrepareLessThan != "-" {
		params.Add("timeToPrepareLessThan", timeToPrepareLessThan)
	}

	if healthyMoreThan != "-" {
		healthyMoreThanInt, err := strconv.Atoi(healthyMoreThan)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value."))
			return
		}
		if healthyMoreThanInt < 1 || healthyMoreThanInt > 3 {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value (Healthy only 1, 2, 3)."))
			return
		}
		params.Add("healthyMoreThan", healthyMoreThan)
	}

	if healthyLessThan != "-" {
		healthyLessThanInt, err := strconv.Atoi(healthyLessThan)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value."))
			return
		}
		if healthyLessThanInt < 1 || healthyLessThanInt > 3 {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value (Healthy only 1, 2, 3)."))
			return
		}
		params.Add("healthyLessThan", healthyLessThan)
	}

	req, err := http.NewRequest("GET", baseUrl+"/api/recipes/filter/", nil)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	//check if response is 200
	if resp.StatusCode != 200 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+resp.Status+" ). Please try again."))
		return
	}

	var recipes []models.Recipe
	err = json.NewDecoder(resp.Body).Decode(&recipes)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	var messageText string
	if len(recipes) == 0 {
		messageText = "No recipes found."
	} else {
		for _, recipe := range recipes {
			messageText += fmt.Sprintf("Id: %d\nTitle: %s\nDescription: %s\nCost: %.2f\nTime to Prepare: %d minutes\nHealthy: %d\n\n", recipe.Id, recipe.Title, recipe.Description, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	bot.Send(msg)
}

// delete recipe
func DeleteRecipe(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository, baseUrl string) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	//get args
	args := strings.Split(update.Message.Text, " ")
	if len(args) < 2 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide the id of the recipe you want to delete."))
		return
	}

	id := args[1]

	req, err := http.NewRequest("DELETE", baseUrl+"/api/recipes/"+id, nil)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while deleting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while deleting recipe ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	//check if response is 200
	if resp.StatusCode != 200 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while deleting recipe ( "+resp.Status+" ). Please try again."))
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Recipe deleted successfully."))
}

func GetFilteredUserRecipes(bot *tgbotapi.BotAPI, update tgbotapi.Update, client *http.Client, CRepo *cache.Repository, baseUrl string) {
	token, err := CRepo.GetKey(strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
		return
	}
	if token == "" {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are not authorized. Please sign in."))
		return
	}

	//get args
	args := strings.Split(update.Message.Text, " ")
	fmt.Println(len(args))
	if len(args) < 7 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide your filter in the format  /getFilteredRecipes <costMoreThat> <costLessThat> <timeToPrepareMoreThan> <timeToPrepareLessThan> <healthyMoreThan> <healthyLessThan> . If some fields are not needed, please provide \"-\"."))
		return
	}

	costMoreThat := args[1]
	costLessThat := args[2]
	timeToPrepareMoreThan := args[3]
	timeToPrepareLessThan := args[4]
	healthyMoreThan := args[5]
	healthyLessThan := args[6]
	params := url.Values{}

	//check if args are valid and add them to params
	if costMoreThat != "-" {
		params.Add("costMoreThat", costMoreThat)
	}

	if costLessThat != "-" {
		params.Add("costLessThat", costLessThat)
	}

	if timeToPrepareMoreThan != "-" {
		params.Add("timeToPrepareMoreThan", timeToPrepareMoreThan)
	}

	if timeToPrepareLessThan != "-" {
		params.Add("timeToPrepareLessThan", timeToPrepareLessThan)
	}

	if healthyMoreThan != "-" {
		healthyMoreThanInt, err := strconv.Atoi(healthyMoreThan)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value."))
			return
		}
		if healthyMoreThanInt < 1 || healthyMoreThanInt > 3 {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value (Healthy only 1, 2, 3)."))
			return
		}
		params.Add("healthyMoreThan", healthyMoreThan)
	}

	if healthyLessThan != "-" {
		healthyLessThanInt, err := strconv.Atoi(healthyLessThan)
		if err != nil {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value."))
			return
		}
		if healthyLessThanInt < 1 || healthyLessThanInt > 3 {
			debug.PrintStack()
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid healthy value. Please provide a valid integer value (Healthy only 1, 2, 3)."))
			return
		}
		params.Add("healthyLessThan", healthyLessThan)
	}

	req, err := http.NewRequest("GET", baseUrl+"/api/recipes/userFilter/", nil)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	defer resp.Body.Close()

	//check if response is 200
	if resp.StatusCode != 200 {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+resp.Status+" ). Please try again."))
		return
	}

	var recipes []models.Recipe
	err = json.NewDecoder(resp.Body).Decode(&recipes)
	if err != nil {
		debug.PrintStack()
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error while getting recipes ( "+err.Error()+" ). Please try again."))
		return
	}

	var messageText string
	if len(recipes) == 0 {
		messageText = "No recipes found."
	} else {
		for _, recipe := range recipes {
			messageText += fmt.Sprintf("Id: %d\nTitle: %s\nDescription: %s\nCost: %.2f\nTime to Prepare: %d minutes\nHealthy: %d\n\n", recipe.Id, recipe.Title, recipe.Description, recipe.Cost, recipe.TimeToPrepare, recipe.Healthy)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
	bot.Send(msg)
}

/*
type RecipesFilter struct {
	CostMoreThan          *float64 `json:"costMoreThan,omitempty"`
	CostLessThan          *float64 `json:"costLessThan,omitempty"`
	TimeToPrepareMoreThan *int     `json:"timeToPrepareMoreThan,omitempty"`
	TimeToPrepareLessThan *int     `json:"timeToPrepareLessThan,omitempty"`
	HealthyMoreThan       *int     `json:"healthyMoreThan,omitempty"`
	HealthyLessThan       *int     `json:"healthyLessThan,omitempty"`
}
		recipes := api.Group("/recipes")
		{
			recipes.POST("/", h.createRecipe)
			recipes.GET("/", h.getAllRecipes)
			recipes.GET("/:id", h.getRecipeById)
			recipes.PUT("/:id", h.updateRecipe)
			recipes.DELETE("/:id", h.deleteRecipe)
			recipes.GET("/public", h.getPublicRecipes)
			recipes.GET("/filter", h.getFilteredRecipes)
			recipes.GET("/userFilter", h.getFilteredUserRecipes)
		}

func (r *RecipesPostgres) GetFilteredRecipes(input models.RecipesFilter) ([]models.Recipe, error) {
	db := r.db
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	var output []models.Recipe

	if input.CostMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost" > $%d`, argId))
		args = append(args, *input.CostMoreThan)
		argId++
	}

	if input.CostLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"cost" < $%d`, argId))
		args = append(args, *input.CostLessThan)
		argId++
	}

	if input.TimeToPrepareMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare" > $%d`, argId))
		args = append(args, *input.TimeToPrepareMoreThan)
		argId++
	}

	if input.TimeToPrepareLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"timeToPrepare" < $%d`, argId))
		args = append(args, *input.TimeToPrepareLessThan)
		argId++
	}

	if input.HealthyMoreThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy" > $%d`, argId))
		args = append(args, *input.HealthyMoreThan)
		argId++
	}

	if input.HealthyLessThan != nil {
		setValues = append(setValues, fmt.Sprintf(`"healthy" < $%d`, argId))
		args = append(args, *input.HealthyLessThan)
		argId++
	}

	setQuery := strings.Join(setValues, " and ")
	if len(setQuery) > 0 {
		setQuery = "and " + setQuery
	}
	log.Println(setQuery)
	log.Println(args...)

	query := fmt.Sprintf(`SELECT rt."id", rt."title", rt."description", rt."public", rt."cost", rt."timeToPrepare", rt."healthy" FROM  %s as rt  WHERE rt.public=true %s`, recipeTable, setQuery)
	args = append(args)

	err := db.Select(&output, query, args...)
	if err != nil {
		return nil, err
	}

	return output, err
}
*/
