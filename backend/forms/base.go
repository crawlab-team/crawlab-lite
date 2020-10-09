package forms

import (
	uuid "github.com/satori/go.uuid"
)

type BaseForm struct {
}

type PageForm struct {
	BaseForm

	PageNum  int `form:"page_num" json:"page_num"`
	PageSize int `form:"page_size" json:"page_size"`
}

func (page *PageForm) Range() (start int, end int) {
	if page.PageNum <= 0 || page.PageSize <= 0 {
		return 0, 0
	}
	start = (page.PageNum - 1) * page.PageSize
	return start, start + page.PageSize - 1
}

type BatchForm struct {
	BaseForm

	Ids []uuid.UUID `form:"ids" json:"ids"`
}
