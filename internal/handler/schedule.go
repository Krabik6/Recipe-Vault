package handler

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) fillSchedule(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.Meal
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Schedule.FillSchedule(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// TODO return id
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) createMeal(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.Meal
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Schedule.CreateMeal(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// TODO return id
	c.JSON(http.StatusOK, gin.H{"id": id})
}

//func (s *ScheduleService) CreateMeal(userId int, meal models.Meal) (int, error) {

func (h *Handler) getAllSchedule(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	output, err := h.services.Schedule.GetAllSchedule(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) getScheduleByPeriod(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	date, ok := c.GetQuery("date")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "date not found")
		return
	}

	period, ok := c.GetQuery("period")
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "period not found")
		return
	}

	intPeriod, err := strconv.Atoi(period)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	output, err := h.services.Schedule.GetScheduleByPeriod(userId, date, intPeriod)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) updateSchedule(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	date, _ := c.GetQuery("date")

	var input models.UpdateScheduleInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Schedule.UpdateSchedule(userId, date, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}

func (h *Handler) deleteSchedule(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	date, _ := c.GetQuery("date")

	err = h.services.Schedule.DeleteSchedule(userId, date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}
