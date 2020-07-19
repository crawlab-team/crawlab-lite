package managers

import (
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	. "github.com/ahmetb/go-linq"
	"github.com/apex/log"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// 任务执行器
type Executor struct {
	Cron *cron.Cron
}

var Exec *Executor

// 初始化
func InitTaskExecutor() error {
	// 构造任务执行器
	c := cron.New(cron.WithSeconds())
	Exec = &Executor{
		Cron: c,
	}

	// 运行定时任务
	if err := Exec.Start(); err != nil {
		return err
	}
	return nil
}

// 任务执行锁
var lockList sync.Map

// 启动任务执行器
func (ex *Executor) Start() error {
	// 启动cron服务
	ex.Cron.Start()

	// 加入执行器到定时任务
	spec := "0/1 * * * * *" // 每秒执行一次
	for i := 0; i < viper.GetInt("task.workers"); i++ {
		// WorkerID
		id := i

		// 初始化任务锁
		lockList.Store(id, false)

		// 加入定时任务
		_, err := ex.Cron.AddFunc(spec, func() {
			ex.ExecuteTask(id)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// 执行任务
func (ex *Executor) ExecuteTask(id int) {
	if flag, ok := lockList.Load(id); ok {
		if flag.(bool) {
			log.Debugf(getWorkerPrefix(id) + "running tasks...")
			return
		}
	}

	// 上锁
	lockList.Store(id, true)

	// 解锁（延迟执行）
	defer func() {
		lockList.Delete(id)
		lockList.Store(id, false)
	}()

	// 开始计时
	tic := time.Now()

	// 获取任务
	task, err := popPendingTask()
	if err != nil {
		log.Errorf("execute task, query task error: %s", err.Error())
		return
	}
	if task == nil {
		return
	}

	// 获取爬虫版本
	var version *models.SpiderVersion
	if task.SpiderVersionId != uuid.Nil {
		version, err = querySpiderVersionById(task.SpiderId, task.SpiderVersionId)
		if err != nil {
			log.Errorf("execute task, query spider version error: %s", err.Error())
			return
		}
	} else {
		version, err = queryLatestSpiderVersion(task.SpiderId)
		if err != nil {
			log.Errorf("execute task, query spider version error: %s", err.Error())
			return
		}
	}
	if version == nil {
		log.Errorf("execute task, query spider version error: spider have no version")
		return
	}

	// 工作目录
	cwd := filepath.Join(
		viper.GetString("spider.path"),
		version.Path,
	)

	// 文件检查
	//if err := SpiderFileCheck(t, spider); err != nil {
	//	log.Errorf("spider file check error: %s", err.Error())
	//	return
	//}

	// 开始执行任务
	log.Infof(getWorkerPrefix(id) + "start task (id:" + task.Id.String() + ")")

	// 更新任务
	task.StartTs = time.Now()                                     // 任务开始时间
	task.Status = constants.TaskStatusRunning                     // 任务状态
	task.WaitDuration = task.StartTs.Sub(task.CreateTs).Seconds() // 等待时长
	if err := updateTask(task); err != nil {
		log.Errorf("execute task, save task error: %s", err.Error())
		return
	}

	// 执行Shell命令
	if err := executeShellCmd(cwd, task); err != nil {
		log.Errorf(getWorkerPrefix(id) + err.Error())

		// 如果发生错误，则发送通知
		//SendNotifications(task, spider)

		return
	}

	// 更新任务
	task.Status = constants.TaskStatusFinished                       // 任务状态: 已完成
	task.FinishTs = time.Now()                                       // 结束时间
	task.RuntimeDuration = task.FinishTs.Sub(task.StartTs).Seconds() // 运行时长
	task.TotalDuration = task.FinishTs.Sub(task.CreateTs).Seconds()  // 总时长
	if err := updateTask(task); err != nil {
		log.Errorf("execute task, save task error: %s", err.Error())
		return
	}

	go finishUpTask(task)

	// 结束计时
	toc := time.Now()

	// 统计时长
	duration := toc.Sub(tic).Seconds()
	durationStr := strconv.FormatFloat(duration, 'f', 6, 64)
	log.Infof(getWorkerPrefix(id) + "task (id:" + task.Id.String() + ")" + " finished. elapsed:" + durationStr + " sec")
}

func StopTask(task *models.Task) {
	if task.Status != constants.TaskStatusRunning &&
		task.Status != constants.TaskStatusProcessing &&
		task.Status != constants.TaskStatusPending {
		return
	}
	ch := utils.TaskExecChanMap.ChanBlocked(task.Id.String())
	if ch != nil {
		ch <- constants.TaskCancel
	}
	task.Status = constants.TaskStatusCancelled                      // 任务状态: 已取消
	task.FinishTs = time.Now()                                       // 结束时间
	task.RuntimeDuration = task.FinishTs.Sub(task.StartTs).Seconds() // 运行时长
	task.TotalDuration = task.FinishTs.Sub(task.CreateTs).Seconds()  // 总时长
}

func getWorkerPrefix(id int) string {
	return "[Worker " + strconv.Itoa(id) + "] "
}

// 执行shell命令
func executeShellCmd(cwd string, task *models.Task) (err error) {
	cmdStr := task.Cmd

	log.Infof("cwd: %s", cwd)
	log.Infof("cmd: %s", cmdStr)

	// 生成执行命令
	var cmd *exec.Cmd
	if runtime.GOOS == constants.Windows {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	// 工作目录
	cmd.Dir = cwd

	// 起一个 goroutine 来监控进程
	ch := utils.TaskExecChanMap.ChanBlocked(task.Id.String())

	// 记录任务生成的日志
	recordTaskLog(cmd, task)

	// 设置环境变量
	setEnv(cmd, task)

	// 完成或取消任务后的工作
	go finishOrCancelTask(ch, cmd, task)

	// kill 的时候，可以 kill 所有的子进程
	if runtime.GOOS != constants.Windows {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// 启动进程
	if err := startTask(cmd, task); err != nil {
		return err
	}

	// 同步等待进程完成
	if err := waitTask(cmd, task); err != nil {
		log.Errorf("task process error: %s", err.Error())
		return err
	}
	ch <- constants.TaskFinish
	return nil
}

func querySpiderVersionById(spiderId uuid.UUID, versionId uuid.UUID) (version *models.SpiderVersion, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		if version, err = tx.SelectSpiderVersion(spiderId, versionId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return version, nil
}

func queryLatestSpiderVersion(spiderId uuid.UUID) (version *models.SpiderVersion, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		versions, err := tx.SelectAllSpiderVersions(spiderId)
		if err != nil {
			return err
		}
		versionI := From(versions).OrderByDescendingT(func(version *models.SpiderVersion) int64 {
			return version.CreateTs.UnixNano()
		}).First()
		if versionI != nil {
			version = versionI.(*models.SpiderVersion)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return version, nil
}

// 设置环境变量
func setEnv(cmd *exec.Cmd, task *models.Task) {
	// 默认把 Node.js 的全局 node_modules 加入环境变量
	envPath := os.Getenv("PATH")
	homePath := os.Getenv("HOME")
	nodeVersion := "v10.19.0"
	nodePath := path.Join(homePath, ".nvm/versions/node", nodeVersion, "lib/node_modules")
	if !strings.Contains(envPath, nodePath) {
		_ = os.Setenv("PATH", nodePath+":"+envPath)
	}
	_ = os.Setenv("NODE_PATH", nodePath)

	// 默认环境变量
	cmd.Env = append(os.Environ(), "CRAWLABLITE_TASK_ID="+task.Id.String())
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=0")
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "TZ=Asia/Shanghai")

	//任务环境变量
	//for _, env := range envs {
	//	cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
	//}

	// 全局环境变量
	//variables := models.GetVariableList()
	//for _, variable := range variables {
	//	cmd.Env = append(cmd.Env, variable.Key+"="+variable.Value)
	//}
}
