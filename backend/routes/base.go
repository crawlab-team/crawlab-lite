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
				authGroup.GET("/spiders", controllers.GetSpiderList)                                    // 爬虫列表
				authGroup.POST("/spiders", controllers.CreateSpider)                                    // 创建爬虫
				authGroup.GET("/spiders/:name", controllers.GetSpider)                                  // 爬虫详情
				authGroup.DELETE("/spiders/:name", controllers.DeleteSpider)                            // 删除爬虫
				authGroup.GET("/spiders/:name/versions", controllers.GetSpiderVersionList)              // 爬虫版本列表
				authGroup.POST("/spiders/:name/versions", controllers.UploadSpiderVersion)              // 上传爬虫版本
				authGroup.DELETE("/spiders/:name/versions/:versionId", controllers.DeleteSpiderVersion) // 删除爬虫版本
			}
			// 任务
			{
				authGroup.GET("/tasks", controllers.GetTaskList)            // 任务列表
				authGroup.GET("/tasks/:id", controllers.GetTask)            // 任务详情
				authGroup.POST("/tasks", controllers.CreateTask)            // 创建任务
				authGroup.POST("/tasks/:id/cancel", controllers.CancelTask) // 取消任务
			}
		}
	}
}
