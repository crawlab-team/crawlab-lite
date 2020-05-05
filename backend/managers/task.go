package managers

import (
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/models"
	"github.com/apex/log"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

// 完成或取消任务
func finishOrCancelTask(ch chan string, cmd *exec.Cmd, task models.Task) {
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

	go finishUpTask(task)
}

// 完成任务的收尾工作
func finishUpTask(task models.Task) {
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

func popPendingTask() (task *models.Task, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		if task, err = tx.SelectTaskWhereStatus(constants.TaskStatusPending); err != nil {
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
