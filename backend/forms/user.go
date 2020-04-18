package forms

type UserForm struct {
	Username string `form:"username" json:"username" binding:"required,min=4,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=5"`
}
