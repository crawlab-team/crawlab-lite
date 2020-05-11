package forms

import uuid "github.com/satori/go.uuid"

type TaskForm struct {
	BaseForm

	SpiderId        uuid.UUID `form:"spider_id" json:"spider_id" binding:"required"`
	SpiderVersionId uuid.UUID `form:"spider_version_id" json:"spider_version_id"`
	ScheduleId      uuid.UUID `form:"spider_version_id" json:"spider_version_id"`
	Cmd             string    `form:"cmd" json:"cmd" binding:"required,min=1"`
}
