package forms

import "mime/multipart"

type ProjectForm struct {
	BaseForm

	Name string `form:"name" json:"name" binding:"required,min=1,max=32"`
}

type ProjectUploadForm struct {
	BaseForm

	File *multipart.FileHeader `form:"file" json:"-" binding:"required"`
}
