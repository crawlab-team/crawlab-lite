package routes

import (
	"crawlab-lite/controllers"
	"crawlab-lite/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	app.Use(middlewares.CORSMiddleware())

	root := app.Group("/api")
	{
		anonymousGroup := root.Group("/")
		{
			anonymousGroup.POST("/login", controllers.Login)                     // 用户登录
			anonymousGroup.GET("/setting", controllers.GetSetting)               // 获取配置信息
			anonymousGroup.GET("/version", controllers.GetVersion)               // 获取发布的版本
			anonymousGroup.GET("/releases/latest", controllers.GetLatestRelease) // 获取最近发布的版本
		}

		authGroup := root.Group("/", middlewares.AuthorizationMiddleware())
		{
			// 爬虫
			{
				authGroup.GET("/spiders", controllers.QuerySpiderList)       // 爬虫列表
				authGroup.GET("/spiders/:name", controllers.QuerySpider)     // 爬虫详情
				authGroup.POST("/spiders", controllers.UploadSpider)         // 上传爬虫
				authGroup.DELETE("/spiders/:name", controllers.DeleteSpider) // 删除爬虫
			}
		}
	}
}
