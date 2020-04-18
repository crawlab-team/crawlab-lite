package routes

import (
	"crawlab-lite/controllers"
	"crawlab-lite/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	app.Use(middlewares.CORSMiddleware())
	anonymousGroup := app.Group("/api")
	{
		anonymousGroup.POST("/login", controllers.Login)                     // 用户登录
		anonymousGroup.GET("/setting", controllers.GetSetting)               // 获取配置信息
		anonymousGroup.GET("/version", controllers.GetVersion)               // 获取发布的版本
		anonymousGroup.GET("/releases/latest", controllers.GetLatestRelease) // 获取最近发布的版本
	}
}
