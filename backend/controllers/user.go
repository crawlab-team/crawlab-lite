package controllers

import (
	"crawlab-lite/constants"
	"crawlab-lite/forms"
	"crawlab-lite/services"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func Login(c *gin.Context) {
	// 绑定请求数据
	var form forms.UserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 校验用户名与密码
	ok, err := services.CheckUser(form.Username, form.Password)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	if ok == false {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	// 生成 Token
	tokenStr, err := services.MakeToken(form.Username)
	if err != nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
		return
	}

	HandleSuccess(c, tokenStr)
}

func GetMe(c *gin.Context) {
	user, exists := c.Get(constants.ContextUser)
	if exists == false || user == nil {
		HandleError(http.StatusUnauthorized, c, errors.New("not authorized"))
	}
	HandleSuccess(c, user)
}
