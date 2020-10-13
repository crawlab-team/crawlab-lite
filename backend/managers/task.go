package managers

import (
	"bufio"
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/database"
	"crawlab-lite/managers/sys_exec"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"errors"
	"github.com/apex/log"
	"os/exec"
	"time"
)

func CancelRunningTasks() (err error) {
	if err := dao.WriteTx(database.MainDB, func(tx dao.Tx) error {
		tasks, err := tx.SelectAllTasks()
		if err != nil {
			return err
		}
		for _, task := range tasks {
			if task.Status == constants.TaskStatusRunning {
				setTaskCancelled(task)
				if err = tx.UpdateTask(task); err != nil {
					return err
				}
			} else if task.Status == constants.TaskStatusProcessing {
				task.Status = constants.TaskStatusPending
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
		task.Error = "start task error: " + err.Error()
		task.Status = constants.TaskStatusError
		task.FinishTs = time.Now()
		_ = updateTask(task)
		return errors.New(task.Error)
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

func finishOrCancelTask(ch chan string, cmd *exec.Cmd, task *models.Task, workerId int) {
	// 传入信号，此处阻塞
	signal := <-ch
	log.Debugf(logPrefix(workerId, task.Id)+"process received signal: %s", signal)

	if signal == constants.TaskCancel && cmd.Process != nil {
		// 终止进程
		if err := sys_exec.KillProcess(cmd); err != nil {
			log.Errorf(logPrefix(workerId, task.Id)+"process kill error: %s", err.Error())

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
	if err := dao.WriteTx(database.MainDB, func(tx dao.Tx) error {
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
	if err := dao.WriteTx(database.MainDB, func(tx dao.Tx) error {
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

func recordTaskLog(cmd *exec.Cmd, task *models.Task) error {
	// get stdout reader
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New("get stdout error: " + err.Error())
	}
	stdoutScanner := bufio.NewScanner(stdout)

	// get stderr scanner
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.New("get stdout error: " + err.Error())
	}
	stderrScanner := bufio.NewScanner(stderr)

	// scan stdout and record
	go func() {
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			if err := saveTaskLog(&models.TaskLog{
				TaskId:   task.Id,
				LineText: text,
				LogStd:   constants.TaskLogStdOut,
				CreateTs: time.Now(),
			}); err != nil {
				break
			}
		}
	}()

	// scan stderr and record
	go func() {
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			if err := saveTaskLog(&models.TaskLog{
				TaskId:   task.Id,
				LineText: text,
				LogStd:   constants.TaskLogStdErr,
				CreateTs: time.Now(),
			}); err != nil {
				break
			}
		}
	}()

	return nil
}

// 保存日志
func saveTaskLog(taskLog *models.TaskLog) (err error) {
	if err := dao.WriteTx(database.LogDB, func(tx dao.Tx) error {
		if err := tx.InsertTaskLog(taskLog); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// 将任务设置为取消状态，并且发送取消信号
func SetTaskCancelled(task *models.Task) {
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

// 将任务设置为取消状态
func setTaskCancelled(task *models.Task) {
	if task.Status != constants.TaskStatusRunning &&
		task.Status != constants.TaskStatusProcessing &&
		task.Status != constants.TaskStatusPending {
		return
	}
	task.Status = constants.TaskStatusCancelled                      // 任务状态: 已取消
	task.FinishTs = time.Now()                                       // 结束时间
	task.RuntimeDuration = task.FinishTs.Sub(task.StartTs).Seconds() // 运行时长
	task.TotalDuration = task.FinishTs.Sub(task.CreateTs).Seconds()  // 总时长
}

func removeTasksOlderThan(days int) (err error) {
	if err := dao.WriteTx(database.MainDB, func(tx dao.Tx) error {
		tasks, err := tx.SelectAllTasks()
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			return nil
		}
		if err := dao.WriteTx(database.LogDB, func(tx2 dao.Tx) error {
			expireDate := time.Now().AddDate(0, 0, -days)
			for _, task := range tasks {
				if task.CreateTs.Before(expireDate) {
					// 取消任务
					SetTaskCancelled(task)

					// 删除任务
					if err = tx.DeleteTask(task.Id); err != nil {
						return err
					}

					// 删除日志
					if err = tx2.DeleteAllTaskLogs(task.Id); err != nil {
						return err
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func removeTaskLogsOlderThan(days int) (err error) {
	if err := dao.ReadTx(database.MainDB, func(tx dao.Tx) error {
		tasks, err := tx.SelectAllTasks()
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			return nil
		}
		if err := dao.WriteTx(database.LogDB, func(tx2 dao.Tx) error {
			for _, task := range tasks {
				if err := tx2.DeleteTaskLogsOlderThan(task.Id, days); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
