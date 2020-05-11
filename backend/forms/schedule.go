package forms

import uuid "github.com/satori/go.uuid"

type ScheduleCreateForm struct {
	BaseForm

	SpiderId        uuid.UUID `form:"spider_id" json:"spider_id" binding:"required"`
	SpiderVersionId uuid.UUID `form:"spider_version_id" json:"spider_version_id"`
	Cron            string    `form:"cron" json:"cron" binding:"required"`
	Cmd             string    `form:"cmd" json:"cmd" binding:"required"`
	Description     string    `form:"description" json:"description"`
}

type ScheduleUpdateForm struct {
	BaseForm

	Cron        string `form:"cron" json:"cron"`
	Cmd         string `form:"cmd" json:"cmd"`
	Enabled     bool   `form:"enabled" json:"enabled"`
	Description string `form:"description" json:"description"`
}
