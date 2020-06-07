package forms

type BaseForm struct {
}

type PageForm struct {
	BaseForm

	PageNum  int `form:"page_num,default=1" json:"page_num" binding:"min=1"`
	PageSize int `form:"page_size,default=10" json:"page_size" binding:"min=1,max=100"`
}

func (page *PageForm) Range() (start int, end int) {
	start = (page.PageNum - 1) * page.PageSize
	return start, start + page.PageSize - 1
}
