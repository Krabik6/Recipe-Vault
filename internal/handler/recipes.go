package handler

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func (h *Handler) createRecipe(c *gin.Context) {
	var input models.Recipe
	if err := c.BindJSON(&input); err != nil {
		log.Print(err)
		return
	}

	err := h.services.Recipes.CreateRecipe(1, input)
	if err != nil {
		log.Print(err)
		return
	}

	c.JSON(200, gin.H{"handler": "createRecipe"})

}

func (h *Handler) getRecipeById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print(err)
		return
	}

	output, err := h.services.Recipes.GetRecipeById(1, id)
	if err != nil {
		log.Print(err)
		return
	}
	c.JSON(200, gin.H{"output": output})

}

func (h *Handler) getAllRecipes(c *gin.Context) {
	output, err := h.services.Recipes.GetAllRecipes(1)
	if err != nil {
		log.Print(err)
		return
	}
	c.JSON(200, gin.H{"output": output})
}

func (h *Handler) updateRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print(err)
		return
	}
	var input models.UpdateRecipeInput
	if err := c.BindJSON(&input); err != nil {
		log.Print(err)
		return
	}

	err = h.services.Recipes.UpdateRecipe(1, id, input)
	if err != nil {
		log.Print(err)
		return
	}
	c.JSON(200, gin.H{"updateRecipe": id})
}

func (h *Handler) deleteRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Print(err)
		return
	}

	err = h.services.Recipes.DeleteRecipe(1, id)
	if err != nil {
		log.Print(err)
		return
	}
	c.JSON(200, gin.H{"deleteRecipe": id})
}
