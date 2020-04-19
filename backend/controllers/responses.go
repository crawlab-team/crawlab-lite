package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ListResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
	Message string      `json:"message"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Data:    data,
		Message: "success",
	})
}

func HandleSuccessList(ctx *gin.Context, total int, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, ListResponse{
		Code:    http.StatusOK,
		Data:    data,
		Total:   total,
		Message: "success",
	})
}

func HandleError(code int, ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(code, Response{
		Code:    code,
		Message: err.Error(),
	})
}
