package routes

import (
	"crawlab-lite/controllers"
	"crawlab-lite/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	app.Use(middlewares.CORSMiddleware())

	root := app.Group("/")
	{
		anonymousGroup := root.Group("/")
		{
			anonymousGroup.POST("/login", controllers.Login)                     // 用户登录
			anonymousGroup.GET("/version", controllers.GetVersion)               // 获取发布的版本
			anonymousGroup.GET("/releases/latest", controllers.GetLatestRelease) // 获取最近发布的版本
		}

		authGroup := root.Group("/", middlewares.AuthorizationMiddleware())
		{
			// 用户
			{
				authGroup.GET("/me", controllers.GetMe) // 获取自己账户
			}
			// 爬虫
			{
				authGroup.GET("/spiders", controllers.GetSpiderList)                                  // 爬虫列表
				authGroup.POST("/spiders", controllers.CreateSpider)                                  // 创建爬虫
				authGroup.GET("/spiders/:id", controllers.GetSpider)                                  // 爬虫详情
				authGroup.DELETE("/spiders/:id", controllers.DeleteSpider)                            // 删除爬虫
				authGroup.GET("/spiders/:id/versions", controllers.GetSpiderVersionList)              // 爬虫版本列表
				authGroup.POST("/spiders/:id/versions", controllers.CreateSpiderVersion)              // 上传爬虫版本
				authGroup.DELETE("/spiders/:id/versions/:versionId", controllers.DeleteSpiderVersion) // 删除爬虫版本
			}
			// 任务
			{
				authGroup.GET("/tasks", controllers.GetTaskList)                  // 任务列表
				authGroup.GET("/tasks/:id", controllers.GetTask)                  // 任务详情
				authGroup.POST("/tasks", controllers.CreateTask)                  // 创建任务
				authGroup.DELETE("/tasks/:id", controllers.DeleteTask)            // 删除任务
				authGroup.POST("/tasks/:id/cancel", controllers.UpdateTaskCancel) // 取消任务
				authGroup.POST("/tasks/:id/restart", controllers.PostTaskRestart) // 重启任务
				authGroup.GET("/tasks/:id/logs", controllers.GetTaskLogList)      // 获取任务日志
			}
			// 定时调度
			{
				authGroup.GET("/schedules", controllers.GetScheduleList)       // 调度列表
				authGroup.GET("/schedules/:id", controllers.GetSchedule)       // 调度详情
				authGroup.POST("/schedules", controllers.CreateSchedule)       // 创建调度
				authGroup.PUT("/schedules/:id", controllers.UpdateSchedule)    // 更新调度
				authGroup.DELETE("/schedules/:id", controllers.DeleteSchedule) // 删除调度
			}
		}
	}
}
