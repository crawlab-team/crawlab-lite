package task

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"github.com/apex/log"
	"os/exec"
	"time"
)

func StartTaskProcess(cmd *exec.Cmd, task models.Task) error {
	if err := cmd.Start(); err != nil {
		log.Errorf("start spider error:{}", err.Error())

		task.Error = "start task error: " + err.Error()
		task.Status = constants.TaskStatusError
		task.FinishTs = time.Now()
		_ = updateTask(&task)
		return err
	}
	return nil
}

func WaitTaskProcess(cmd *exec.Cmd, task models.Task) error {
	if err := cmd.Wait(); err != nil {
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
				_ = updateTask(&task)

				go FinishUpTask(task)
			}
		}

		return err
	}

	return nil
}
