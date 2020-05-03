package forms

type TaskForm struct {
	BaseForm

	SpiderName string `form:"spider_name" json:"spider_name" binding:"required,min=1,max=32"`
	Cmd        string `form:"cmd" json:"cmd" binding:"required,min=1"`
}
