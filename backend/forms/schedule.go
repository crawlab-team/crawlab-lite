package forms

type ScheduleCreateForm struct {
	BaseForm

	SpiderName      string `form:"spider_name" json:"spider_name" binding:"required,min=1,max=32"`
	SpiderVersionId string `form:"spider_version_id" json:"spider_version_id"`
	Cron            string `form:"cron" json:"cron" binding:"required,min=1"`
	Cmd             string `form:"cmd" json:"cmd" binding:"required,min=1"`
	Description     string `form:"description" json:"description"`
}

type ScheduleUpdateForm struct {
	BaseForm

	Cron        string `form:"cron" json:"cron"`
	Cmd         string `form:"cmd" json:"cmd"`
	Enabled     bool   `form:"enabled" json:"enabled"`
	Description string `form:"description" json:"description"`
}
