package forms

type BaseForm struct {
}

type PageForm struct {
	BaseForm

	PageNum  int `form:"pagenum,default=1" json:"pagenum" binding:"min=1"`
	PageSize int `form:"pagesize,default=10" json:"pagesize" binding:"min=1,max=100"`
}
