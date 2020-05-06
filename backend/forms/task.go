package forms

import uuid "github.com/satori/go.uuid"

type TaskForm struct {
	BaseForm

	SpiderName      string    `form:"spider_name" json:"spider_name" binding:"required,min=1,max=32"`
	SpiderVersionId string    `form:"spider_version_id" json:"spider_version_id"`
	ScheduleId      uuid.UUID `form:"spider_version_id" json:"spider_version_id"`
	Cmd             string    `form:"cmd" json:"cmd" binding:"required,min=1"`
}
