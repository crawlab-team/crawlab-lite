package main

import (
	"context"
	"crawlab-lite/config"
	"crawlab-lite/database"
	"crawlab-lite/lib/validate_bridge"
	"crawlab-lite/routes"
	"crawlab-lite/task"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	binding.Validator = new(validate_bridge.DefaultValidator)
	app := gin.Default()

	// 初始化配置
	if err := config.InitConfig("./config.yml"); err != nil {
		log.Error("Init config error:" + err.Error())
		panic(err)
	}
	log.Info("Initialized config successfully")

	// 初始化 Key-Value 数据库
	if err := database.InitKvDB(); err != nil {
		log.Error("Init key-value database error:" + err.Error())
		panic(err)
	}
	log.Info("Initialized key-value database successfully")

	// 初始化日志设置
	logLevel := viper.GetString("log.level")
	if logLevel != "" {
		log.SetLevelFromString(logLevel)
	}
	log.Info("Initialized log config successfully")

	// 初始化任务执行器
	if err := task.InitTaskExecutor(); err != nil {
		log.Error("Init task executor error:" + err.Error())
		panic(err)
	}
	log.Info("Initialized task executor successfully")

	// 初始化路由
	routes.InitRoutes(app)

	// 运行服务器
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)
	srv := &http.Server{
		Handler: app,
		Addr:    address,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error("Run server error:" + err.Error())
			} else {
				log.Info("Server graceful down")
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Error("Run server error:" + err.Error())
	}
}
