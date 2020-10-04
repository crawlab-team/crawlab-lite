package managers

import (
	"errors"
	"github.com/apex/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

// 数据清扫器
var Clearer *clearer

type clearer struct {
	cron *cron.Cron
}

// 启动数据清扫器
func (cl *clearer) Start() error {
	// 先执行一次清扫
	if err := cl.Clear(); err != nil {
		return err
	}

	// 启动cron服务
	cl.cron.Start()

	// 每小时一次周期执行清扫
	if _, err := cl.cron.AddFunc("0 */60 * * * *", func() {
		go func() {
			if err := cl.Clear(); err != nil {
				log.Error(err.Error())
			}
		}()
	}); err != nil {
		return err
	}

	return nil
}

// 执行清扫
func (cl *clearer) Clear() error {
	// 清理过期的任务
	log.Infof("clearing expired task...")
	taskDays := viper.GetInt("task.expireDays")
	if err := removeTasksOlderThan(taskDays); err != nil {
		return errors.New("clear expired tasks error: " + err.Error())
	}

	// 清理过期的任务日志
	log.Infof("clearing expired task logs...")
	logDays := viper.GetInt("log.expireDays")
	if err := removeTaskLogsOlderThan(logDays); err != nil {
		return errors.New("clear expired task logs error: " + err.Error())
	}
	return nil
}

// 初始化数据清扫器
func InitClearer() error {
	Clearer = &clearer{
		cron: cron.New(cron.WithSeconds()),
	}
	if err := Clearer.Start(); err != nil {
		return err
	}
	return nil
}
