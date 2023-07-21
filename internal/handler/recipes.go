package handler

import (
	"encoding/json"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (h *Handler) createRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the uploaded file(s)
	files := c.Request.MultipartForm.File["images"]

	// Create an empty slice to store the image file headers
	imageFiles := make([]*multipart.FileHeader, len(files))

	// Process each file
	for i, file := range files {
		imageFiles[i] = file
	}

	// Получите строку JSON из формы
	ingredientsJSON := c.PostForm("ingredients")

	// Создайте пустой массив для ингредиентов
	var ingredientInputs []models.IngredientInput

	// Декодируйте строку JSON обратно в массив
	err = json.Unmarshal([]byte(ingredientsJSON), &ingredientInputs)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Теперь у вас есть массив ingredientInputs, который вы можете использовать
	recipe := models.RecipeInput{
		Title:            c.PostForm("title"),
		Description:      c.PostForm("description"),
		IsPublic:         c.PostForm("public") == "true",
		Cost:             parseStringToFloat64(c.PostForm("cost")),
		TimeToPrepare:    parseStringToInt64(c.PostForm("timeToPrepare")),
		Healthy:          parseStringToInt64(c.PostForm("healthy")),
		IngredientInputs: ingredientInputs,
	}

	// Create the recipe in the service
	id, err := h.services.Recipes.CreateRecipe(userId, recipe, imageFiles)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// func parse cost that returns float64 and takes string
func parseStringToFloat64(valStr string) float64 {
	valFloat64, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0
	}
	return valFloat64
}

// parse string to int64
func parseStringToInt64(valStr string) int64 {
	valInt64, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0
	}
	return valInt64
}

func (h *Handler) getRecipeById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	output, err := h.services.Recipes.GetRecipeById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)

}

func (h *Handler) getAllRecipes(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	output, err := h.services.Recipes.GetAllRecipes(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) getFilteredUserRecipes(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.RecipesFilter
	costMoreThat := c.Query("costMoreThat")
	if costMoreThat != "" {
		// cost to float64
		cost, err := strconv.ParseFloat(costMoreThat, 64)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.CostMoreThan = &cost
	}

	costLessThat := c.Query("costLessThat")
	if costLessThat != "" {
		cost, err := strconv.ParseFloat(costMoreThat, 64)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.CostLessThan = &cost
	}

	timeToPrepareMoreThan := c.Query("timeToPrepareMoreThan")
	if timeToPrepareMoreThan != "" {
		time, err := strconv.Atoi(timeToPrepareMoreThan)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.TimeToPrepareMoreThan = &time
	}

	timeToPrepareLessThan := c.Query("timeToPrepareLessThan")
	if timeToPrepareLessThan != "" {
		time, err := strconv.Atoi(timeToPrepareLessThan)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.TimeToPrepareLessThan = &time
	}

	output, err := h.services.Recipes.GetFilteredUserRecipes(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) getPublicRecipes(c *gin.Context) {
	output, err := h.services.Recipes.GetPublicRecipes()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) getFilteredRecipes(c *gin.Context) {
	var input models.RecipesFilter
	costMoreThat := c.Query("costMoreThat")
	if costMoreThat != "" {
		// cost to float64
		cost, err := strconv.ParseFloat(costMoreThat, 64)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.CostMoreThan = &cost
	}

	costLessThat := c.Query("costLessThat")
	if costLessThat != "" {
		cost, err := strconv.ParseFloat(costMoreThat, 64)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.CostLessThan = &cost
	}

	timeToPrepareMoreThan := c.Query("timeToPrepareMoreThan")
	if timeToPrepareMoreThan != "" {
		time, err := strconv.Atoi(timeToPrepareMoreThan)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.TimeToPrepareMoreThan = &time
	}

	timeToPrepareLessThan := c.Query("timeToPrepareLessThan")
	if timeToPrepareLessThan != "" {
		time, err := strconv.Atoi(timeToPrepareLessThan)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.TimeToPrepareLessThan = &time
	}

	healthyMoreThan := c.Query("healthyMoreThan")
	if healthyMoreThan != "" {
		healthy, err := strconv.Atoi(healthyMoreThan)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.HealthyMoreThan = &healthy
	}

	healthyLessThan := c.Query("healthyLessThan")
	if healthyLessThan != "" {
		healthy, err := strconv.Atoi(healthyLessThan)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		input.HealthyLessThan = &healthy
	}

	output, err := h.services.Recipes.GetFilteredRecipes(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)

}

func (h *Handler) updateRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get the uploaded file(s)
	files := c.Request.MultipartForm.File["images"]

	// Create an empty slice to store the image file headers
	imageFiles := make([]*multipart.FileHeader, len(files))

	// Process each file
	for i, file := range files {
		imageFiles[i] = file
	}

	// Получите строку JSON из формы
	ingredientsJSON := c.PostForm("ingredients")
	// Создайте пустой массив для ингредиентов
	var ingredientInputs []models.IngredientInput

	// Декодируйте строку JSON обратно в массив
	err = json.Unmarshal([]byte(ingredientsJSON), &ingredientInputs)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Create an instance of the UpdateRecipeInput struct with the input data
	input := models.UpdateRecipeInput{
		Title:            getStringValue(c, "title"),
		Description:      getStringValue(c, "description"),
		IsPublic:         getBoolValue(c, "public"),
		Cost:             getFloat64Value(c, "cost"),
		TimeToPrepare:    getIntValue(c, "timeToPrepare"),
		Healthy:          getIntValue(c, "healthy"),
		IngredientInputs: &ingredientInputs,
	}

	// Update the recipe in the service
	err = h.services.Recipes.UpdateRecipe(userId, id, input, imageFiles)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// Helper functions to retrieve form values
func getStringValue(c *gin.Context, key string) *string {
	value := c.PostForm(key)
	if value == "" {
		return nil
	}
	return &value
}

func getBoolValue(c *gin.Context, key string) *bool {
	value := c.PostForm(key)
	if value == "" {
		return nil
	}
	boolValue, _ := strconv.ParseBool(value)
	return &boolValue
}

func getFloat64Value(c *gin.Context, key string) *float64 {
	value := c.PostForm(key)
	if value == "" {
		return nil
	}
	floatValue, _ := strconv.ParseFloat(value, 64)
	return &floatValue
}

func getIntValue(c *gin.Context, key string) *int {
	value := c.PostForm(key)
	if value == "" {
		return nil
	}
	intValue, _ := strconv.Atoi(value)
	return &intValue
}

func (h *Handler) deleteRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Recipes.DeleteRecipe(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}
