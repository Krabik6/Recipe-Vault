package handler

import (
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) fillSchedule(c *gin.Context) {
	var input models.Schedule
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Schedule.FillSchedule(1, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// TODO return id
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) getAllSchedule(c *gin.Context) {
	output, err := h.services.Schedule.GetAllSchedule(1)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) getScheduleByDate(c *gin.Context) {
	date, _ := c.GetQuery("date")
	log.Println(date)

	output, err := h.services.Schedule.GetScheduleByDate(1, date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *Handler) updateSchedule(c *gin.Context) {
	date, _ := c.GetQuery("date")

	var input models.UpdateScheduleInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Schedule.UpdateSchedule(1, date, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}

func (h *Handler) deleteSchedule(c *gin.Context) {
	date, _ := c.GetQuery("date")

	err := h.services.Schedule.DeleteSchedule(1, date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}
