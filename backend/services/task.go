package services

import (
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/forms"
	"crawlab-lite/models"
	"errors"
	uuid "github.com/satori/go.uuid"
)

func QueryTaskPage(page forms.PageForm) (total int, tasks []*models.Task, err error) {
	start, end := page.Range()

	if err := dao.ReadTx(func(tx dao.Tx) error {
		if tasks, err = tx.SelectAllTasksLimit(start, end); err != nil {
			return err
		}
		if total, err = tx.CountTasks(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, tasks, nil
}

func QueryTaskById(id uuid.UUID) (task *models.Task, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		if task, err = tx.SelectTask(id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return task, nil
}

func AddTask(form forms.TaskForm) (task *models.Task, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查爬虫是否存在
		if spider, err := tx.SelectSpider(form.SpiderId); err != nil {
			return err
		} else if spider == nil {
			return errors.New("spider not found")
		}

		if form.SpiderVersionId != uuid.Nil {
			// 检查爬虫版本是否存在
			version, err := tx.SelectSpiderVersion(form.SpiderId, form.SpiderVersionId)
			if err != nil {
				return err
			} else if version == nil {
				return errors.New("spider version not found")
			}
		}

		task = &models.Task{
			SpiderId:        form.SpiderId,
			SpiderVersionId: form.SpiderVersionId,
			ScheduleId:      form.ScheduleId,
			Cmd:             form.Cmd,
		}

		// 存储任务信息
		if err := tx.InsertTask(task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return task, nil
}

func CancelTask(id uuid.UUID, status constants.TaskStatus) (task *models.Task, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		if task, err = tx.SelectTask(id); err != nil {
			return err
		}
		if task == nil {
			return errors.New("task not found")
		}
		if task.Status == constants.TaskStatusFinished {
			return errors.New("task has been finished")
		}
		task.Status = status
		if err = tx.UpdateTask(task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return task, nil
}
