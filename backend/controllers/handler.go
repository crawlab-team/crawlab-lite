package controllers

import (
	"crawlab-lite/results"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleSuccess(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, results.ResponseBody{
		Code:    http.StatusOK,
		Data:    data,
		Message: "success",
	})
}

func HandleSuccessList(ctx *gin.Context, total int, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, results.ResponseBody{
		Code: http.StatusOK,
		Data: struct {
			List  interface{} `json:"list"`
			Total int         `json:"total"`
		}{
			List:  data,
			Total: total,
		},
		Message: "success",
	})
}

func HandleError(code int, ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(code, results.ResponseBody{
		Code:    code,
		Message: err.Error(),
	})
}
