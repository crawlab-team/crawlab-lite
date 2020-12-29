package controllers

import (
	"crawlab-lite/forms"
	"crawlab-lite/services"
	"errors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func GetScheduleList(c *gin.Context) {
	var page forms.PageForm

	if err := c.ShouldBindQuery(&page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if total, schedules, err := services.QuerySchedulePage(page); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccessList(c, total, schedules)
}

func GetSchedule(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
		return
	}

	if schedule, err := services.QueryScheduleById(id); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if schedule == nil {
		HandleError(http.StatusNotFound, c, errors.New("schedule not found"))
		return
	}

	HandleSuccess(c, schedule)
}

func CreateSchedule(c *gin.Context) {
	var form forms.ScheduleCreateForm

	if err := c.ShouldBindJSON(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if res, err := services.AddSchedule(form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, res)
}

func UpdateSchedule(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
		return
	}

	var form forms.ScheduleUpdateForm

	if err := c.ShouldBindJSON(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if res, err := services.ModifySchedule(id, form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, res)
}

func DeleteSchedule(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
		return
	}

	if res, err := services.RemoveSchedule(id); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	HandleSuccess(c, res)
}
