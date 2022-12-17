package handler

import "github.com/gin-gonic/gin"

func (h *Handler) addRecipe(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "addRecipe"})
}

func (h *Handler) deleteRecipe(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "deleteRecipe"})
}

func (h *Handler) updateRecipe(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "updateRecipe"})
}

func (h *Handler) getAllRecipes(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "getAllRecipes"})
}

func (h *Handler) getRecipeById(c *gin.Context) {
	c.JSON(200, gin.H{"handler": "getRecipeById"})
}
