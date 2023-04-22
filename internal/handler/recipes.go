package handler

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) createRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	var input models.Recipe
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Recipes.CreateRecipe(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
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
	var input models.UpdateRecipeInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Recipes.UpdateRecipe(userId, id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
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
