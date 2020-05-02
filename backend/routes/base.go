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
			// 爬虫项目
			{
				authGroup.GET("/projects", controllers.GetProjectList)                                    // 爬虫项目列表
				authGroup.POST("/projects", controllers.CreateProject)                                    // 创建爬虫项目
				authGroup.GET("/projects/:name", controllers.GetProject)                                  // 爬虫项目详情
				authGroup.DELETE("/projects/:name", controllers.DeleteProject)                            // 删除爬虫项目
				authGroup.GET("/projects/:name/versions", controllers.GetProjectVersionList)              // 项目版本列表
				authGroup.POST("/projects/:name/versions", controllers.UploadProjectVersion)              // 上传项目版本
				authGroup.DELETE("/projects/:name/versions/:versionId", controllers.DeleteProjectVersion) // 删除项目版本
			}
		}
	}
}
