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
	} else {
		HandleSuccessList(c, total, schedules)
	}
}

func GetSchedule(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
	}

	if schedule, err := services.QueryScheduleById(id); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	} else {
		if schedule == nil {
			HandleError(http.StatusNotFound, c, errors.New("schedule not found"))
			return
		}
		HandleSuccess(c, schedule)
	}
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
	} else {
		HandleSuccess(c, res)
	}
}

func UpdateSchedule(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
	}

	var form forms.ScheduleUpdateForm

	if err := c.ShouldBindJSON(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if res, err := services.ModifySchedule(id, form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}

func DeleteSchedule(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		HandleError(http.StatusBadRequest, c, errors.New("invalid id"))
	}

	if res, err := services.RemoveSchedule(id); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	} else {
		HandleSuccess(c, res)
	}
}
