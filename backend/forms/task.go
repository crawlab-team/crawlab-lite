package forms

type TaskForm struct {
	BaseForm

	SpiderName      string `form:"spider_name" json:"spider_name" binding:"required,min=1,max=32"`
	SpiderVersionId string `form:"spider_version_id" json:"spider_version_id"`
	ScheduleId      string `form:"spider_version_id" json:"spider_version_id"`
	Cmd             string `form:"cmd" json:"cmd" binding:"required,min=1"`
}
