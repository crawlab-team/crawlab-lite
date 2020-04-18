package middlewares

import (
	"crawlab-lite/constants"
	"crawlab-lite/controllers"
	"crawlab-lite/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Token string
		tokenStr := c.GetHeader("Authorization")

		// 校验 Token 并从其中提取用户
		user, err := services.GetUserFromToken(tokenStr)

		// 校验失败，返回错误响应
		if err != nil {
			controllers.HandleError(http.StatusUnauthorized, c, errors.New("unauthorized"))
			return
		}

		c.Set(constants.ContextUser, &user)

		// 校验成功
		c.Next()
	}
}
