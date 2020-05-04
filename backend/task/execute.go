package task

import (
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/models"
	"crawlab-lite/services"
	"crawlab-lite/utils"
	"github.com/apex/log"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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
var LockList sync.Map

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
		LockList.Store(id, false)

		// 加入定时任务
		_, err := ex.Cron.AddFunc(spec, func() {
			ExecuteTask(id)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// 执行任务
func ExecuteTask(id int) {
	if flag, ok := LockList.Load(id); ok {
		if flag.(bool) {
			log.Debugf(GetWorkerPrefix(id) + "running tasks...")
			return
		}
	}

	// 上锁
	LockList.Store(id, true)

	// 解锁（延迟执行）
	defer func() {
		LockList.Delete(id)
		LockList.Store(id, false)
	}()

	// 开始计时
	tic := time.Now()

	// 获取任务
	task, err := services.PopPendingTask()
	if err != nil {
		log.Errorf("execute task, query task error: %s", err.Error())
		return
	}
	if task == nil {
		return
	}

	// 获取爬虫版本
	var version *models.SpiderVersion
	if task.SpiderVersionId != "" {
		version, err = services.QuerySpiderVersionById(task.SpiderName, task.SpiderVersionId)
		if err != nil {
			log.Errorf("execute task, query spider version error: %s", err.Error())
			return
		}
	} else {
		version, err = services.QueryLatestSpiderVersion(task.SpiderName)
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
	log.Infof(GetWorkerPrefix(id) + "start task (id:" + task.Id + ")")

	// 更新任务
	task.StartTs = time.Now()                                     // 任务开始时间
	task.Status = constants.TaskStatusRunning                     // 任务状态
	task.WaitDuration = task.StartTs.Sub(task.CreateTs).Seconds() // 等待时长
	if err := updateTask(task); err != nil {
		log.Errorf("execute task, save task error: %s", err.Error())
		return
	}

	// 执行Shell命令
	if err := ExecuteShellCmd(cwd, *task); err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())

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

	go FinishUpTask(*task)

	// 结束计时
	toc := time.Now()

	// 统计时长
	duration := toc.Sub(tic).Seconds()
	durationStr := strconv.FormatFloat(duration, 'f', 6, 64)
	log.Infof(GetWorkerPrefix(id) + "task (id:" + task.Id + ")" + " finished. elapsed:" + durationStr + " sec")
}

func GetWorkerPrefix(id int) string {
	return "[Worker " + strconv.Itoa(id) + "] "
}

// 执行shell命令
func ExecuteShellCmd(cwd string, task models.Task) (err error) {
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
	ch := utils.TaskExecChanMap.ChanBlocked(task.Id)

	go FinishOrCancelTask(ch, cmd, task)

	// kill 的时候，可以 kill 所有的子进程
	if runtime.GOOS != constants.Windows {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	// 启动进程
	if err := StartTaskProcess(cmd, task); err != nil {
		return err
	}

	// 同步等待进程完成
	if err := WaitTaskProcess(cmd, task); err != nil {
		return err
	}
	ch <- constants.TaskFinish
	return nil
}

// 完成或取消任务
func FinishOrCancelTask(ch chan string, cmd *exec.Cmd, task models.Task) {
	// 传入信号，此处阻塞
	signal := <-ch
	log.Infof("process received signal: %s", signal)

	if signal == constants.TaskCancel && cmd.Process != nil {
		var err error
		// 兼容windows
		if runtime.GOOS == constants.Windows {
			err = cmd.Process.Kill()
		} else {
			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		// 取消进程
		if err != nil {
			log.Errorf("process kill error: %s", err.Error())

			task.Error = "kill process error: " + err.Error()
			task.Status = constants.TaskStatusError
		} else {
			task.Error = "user kill the process ..."
			task.Status = constants.TaskStatusCancelled
		}
	} else {
		// 保存任务
		task.Status = constants.TaskStatusFinished
	}

	task.FinishTs = time.Now()
	_ = updateTask(&task)

	go FinishUpTask(task)
}

// 完成任务的收尾工作
func FinishUpTask(task models.Task) {
	// TODO
}

func updateTask(task *models.Task) (err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		if err = tx.UpdateTask(task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil
	}
	return nil
}
