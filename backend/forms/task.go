package forms

import (
	"crawlab-lite/constants"
	uuid "github.com/satori/go.uuid"
)

type TaskPageForm struct {
	PageForm

	SpiderId   string               `form:"spider_id" json:"spider_id"`
	ScheduleId string               `form:"schedule_id" json:"schedule_id"`
	Status     constants.TaskStatus `form:"status" json:"status"`
}

type TaskForm struct {
	BaseForm

	SpiderId        uuid.UUID `form:"spider_id" json:"spider_id" binding:"required"`
	SpiderVersionId uuid.UUID `form:"spider_version_id" json:"spider_version_id"`
	ScheduleId      uuid.UUID `form:"schedule_id" json:"schedule_id"`
	Cmd             string    `form:"cmd" json:"cmd" binding:"required,min=1"`
}
