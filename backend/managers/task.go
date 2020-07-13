package managers

import (
	"bufio"
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"github.com/apex/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lestrrat-go/strftime"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func CancelRunningTasks() (err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		tasks, err := tx.SelectAllTasks()
		if err != nil {
			return err
		}
		for _, task := range tasks {
			if task.Status == constants.TaskStatusRunning {
				task.Status = constants.TaskStatusCancelled
				if err = tx.UpdateTask(task); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func startTask(cmd *exec.Cmd, task *models.Task) (err error) {
	if err = cmd.Start(); err != nil {
		log.Errorf("start spider error:{}", err.Error())

		task.Error = "start task error: " + err.Error()
		task.Status = constants.TaskStatusError
		task.FinishTs = time.Now()
		_ = updateTask(task)
		return err
	}
	return nil
}

func waitTask(cmd *exec.Cmd, task *models.Task) (err error) {
	if err = cmd.Wait(); err != nil {
		log.Errorf("wait process finish error: %s", err.Error())

		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Errorf("exit error, exit code: %d", exitCode)

			// 非kill 的错误类型
			if exitCode != -1 {
				// 非手动kill保存为错误状态
				task.Error = err.Error()
				task.FinishTs = time.Now()
				task.Status = constants.TaskStatusError
				_ = updateTask(task)

				go finishUpTask(task)
			}
		}

		return err
	}

	return nil
}

func finishOrCancelTask(ch chan string, cmd *exec.Cmd, task *models.Task) {
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
	_ = updateTask(task)

	go finishUpTask(task)
}

// 完成任务的收尾工作
func finishUpTask(task *models.Task) {
	// TODO
}

func updateTask(task *models.Task) (err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		if err = tx.UpdateTask(task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func popPendingTask() (task *models.Task, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		if task, err = tx.SelectFirstTaskWhereStatus(constants.TaskStatusPending); err != nil {
			return err
		}
		if task == nil {
			return nil
		}
		task.Status = constants.TaskStatusProcessing
		if err = tx.UpdateTask(task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return task, nil
}

func recordTaskLog(cmd *exec.Cmd, task *models.Task) {
	// 创建日志目录
	fileDir, err := makeLogDir(*task)
	if err != nil {
		return
	}

	days := viper.GetInt("log.expireDays")

	logPattern := filepath.Join(fileDir, task.Id.String()) + "_%Y%m%d.log"
	logPath, err := strftime.Format(logPattern, time.Now())
	if err != nil {
		log.Errorf("new strftime error: %s", err.Error())
		return
	}

	logf, err := rotatelogs.New(
		logPattern,
		rotatelogs.WithMaxAge(time.Duration(days*24)*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Errorf("new rotatelogs error: %s", err.Error())
		return
	}

	task.LogPath = logPath
	_ = updateTask(task)

	// get stdout reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("get stdout error: %s", err.Error())
		return
	}
	stdoutScanner := bufio.NewScanner(stdout)

	// get stderr scanner
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Errorf("get stdout error: %s", err.Error())
		return
	}
	stderrScanner := bufio.NewScanner(stderr)

	// scan stdout and record
	go func() {
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			_, err = logf.Write([]byte(text + "\n"))
			if err != nil {
				break
			}
		}
	}()

	// scan stderr and record
	go func() {
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			_, err = logf.Write([]byte(text + "\n"))
			if err != nil {
				break
			}
		}
	}()
}

// 生成日志目录
func makeLogDir(task models.Task) (fileDir string, err error) {
	// 日志目录
	fileDir = filepath.Join(viper.GetString("log.path"), task.SpiderId.String())

	// 如果日志目录不存在，生成该目录
	if !utils.PathExist(fileDir) {
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			log.Errorf("execute task, make log dir error: %s", err.Error())
			return "", err
		}
	}

	return fileDir, nil
}
