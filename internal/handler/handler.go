package handler

import (
	"github.com/Krabik6/meal-schedule/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: *services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

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
			recipes.GET("/filter", h.getFilteredRecipes)
			recipes.GET("/userFilter", h.getFilteredUserRecipes)
		}
		schedule := api.Group("/schedule")
		{
			schedule.POST("/", h.createMeal)
			schedule.GET("/all", h.getAllSchedule)
			schedule.GET("/", h.getScheduleByPeriod)
			schedule.PUT("/", h.updateSchedule)
			schedule.DELETE("/", h.deleteSchedule)
		}

		return router
	}

}
