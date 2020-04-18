package forms

import "mime/multipart"

type SpiderUploadForm struct {
	File *multipart.FileHeader `form:"file" json:"-" binding:"required"`
	Name string                `form:"name" json:"name" binding:"required,min=1,max=32"`
	Cmd  string                `form:"cmd" json:"cmd"`
}

type SpiderDeleteForm struct {
	Name string `form:"name" json:"name" binding:"required,min=1,max=32"`
}
