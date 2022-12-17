package handler

import (
	"github.com/Krabik6/meal-schedule/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: *services}
}

func signup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "signup",
	})
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.GET("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
	}
	recipes := router.Group("/recipes")
	{
		recipes.POST("/", h.addRecipe)
		recipes.GET("/", h.getAllRecipes)
		recipes.GET("/:id", h.addRecipe)
		recipes.PUT("/:id", h.updateRecipe)
		recipes.DELETE("/:id", h.deleteRecipe)

	}
	return router
}
