package services

import (
	"bufio"
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/forms"
	"crawlab-lite/managers"
	"crawlab-lite/models"
	"crawlab-lite/results"
	"errors"
	. "github.com/ahmetb/go-linq"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"os"
)

func QueryTaskPage(page forms.TaskPageForm) (total int, resultList []*results.Task, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		allTasks, err := tx.SelectAllTasks()
		if err != nil {
			return err
		}

		spiderId := uuid.FromStringOrNil(page.SpiderId)
		scheduleId := uuid.FromStringOrNil(page.ScheduleId)

		query := From(allTasks)
		if spiderId != uuid.Nil {
			query = query.WhereT(func(task *models.Task) bool {
				return task.SpiderId == spiderId
			})
		}
		if scheduleId != uuid.Nil {
			query = query.WhereT(func(task *models.Task) bool {
				return task.ScheduleId == scheduleId
			})
		}
		if page.Status != "" {
			query = query.WhereT(func(task *models.Task) bool {
				return task.Status == page.Status
			})
		}

		query = query.OrderByDescendingT(func(task *models.Task) int64 {
			return task.StartTs.UnixNano()
		}).Query
		total = query.Count()

		if page.PageNum > 0 && page.PageSize > 0 {
			start, end := page.Range()
			query = query.Skip(start).Take(end - start)
		}
		tasks := query.Results()

		cache := map[uuid.UUID]*models.Spider{}
		for _, task := range tasks {
			task := task.(*models.Task)
			var result results.Task
			if err := copier.Copy(&result, task); err != nil {
				return err
			}
			var spider *models.Spider
			spider, exists := cache[task.SpiderId]
			if !exists {
				if spider, err = tx.SelectSpider(task.SpiderId); err != nil {
					return err
				}
				if spider != nil {
					cache[task.SpiderId] = spider
				}
			}
			if spider == nil {
				continue
			}
			result.SpiderName = spider.Name
			resultList = append(resultList, &result)
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, resultList, nil
}

func QueryTaskById(id uuid.UUID) (result *results.Task, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		task, err := tx.SelectTask(id)
		if err != nil {
			return err
		}
		spider, err := tx.SelectSpider(task.SpiderId)
		if err != nil {
			return err
		}
		result = &results.Task{}
		if err := copier.Copy(&result, task); err != nil {
			return err
		}
		result.SpiderName = spider.Name
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func AddTask(form forms.TaskForm) (result *results.Task, err error) {
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

		task := &models.Task{
			SpiderId:        form.SpiderId,
			SpiderVersionId: form.SpiderVersionId,
			ScheduleId:      form.ScheduleId,
			Cmd:             form.Cmd,
		}

		// 存储任务信息
		if err := tx.InsertTask(task); err != nil {
			return err
		}

		result = &results.Task{}
		if err := copier.Copy(&result, task); err != nil {
			return err
		}
		spider, err := tx.SelectSpider(task.SpiderId)
		if err != nil {
			return err
		}
		result.SpiderName = spider.Name
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func RemoveTask(id uuid.UUID) (res interface{}, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查任务是否存在
		if task, err := tx.SelectTask(id); err != nil {
			return err
		} else if task == nil {
			return errors.New("task not found")
		} else {
			managers.StopTask(task)

			// 删除任务
			if err = tx.DeleteTask(id); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func CancelTask(id uuid.UUID, status constants.TaskStatus) (task *models.Task, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		if task, err = tx.SelectTask(id); err != nil {
			return err
		}
		if task == nil {
			return errors.New("task not found")
		}

		if task.Status == constants.TaskStatusRunning ||
			task.Status == constants.TaskStatusProcessing ||
			task.Status == constants.TaskStatusPending {
			managers.StopTask(task)
		} else {
			return nil
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

func RestartTask(id uuid.UUID) (result *results.Task, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		var task *models.Task
		if task, err = tx.SelectTask(id); err != nil {
			return err
		}
		if task == nil {
			return errors.New("task not found")
		}

		managers.StopTask(task)

		newTask := &models.Task{
			SpiderId:        task.SpiderId,
			SpiderVersionId: task.SpiderVersionId,
			ScheduleId:      task.ScheduleId,
			Cmd:             task.Cmd,
		}

		// 存储任务信息
		if err := tx.InsertTask(newTask); err != nil {
			return err
		}

		result = &results.Task{}
		if err := copier.Copy(&result, newTask); err != nil {
			return err
		}
		spider, err := tx.SelectSpider(newTask.SpiderId)
		if err != nil {
			return err
		}
		result.SpiderName = spider.Name
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func QueryTaskLogPage(form forms.TaskLogPageForm) (total int, resultList []*results.TaskLogLine, err error) {
	taskId, err := uuid.FromString(form.TaskId)
	if err != nil {
		return 0, nil, err
	}

	if form.PageSize <= 0 || form.PageNum <= 0 {
		return 0, nil, nil
	}

	var task *models.Task
	if err := dao.ReadTx(func(tx dao.Tx) error {
		if task, err = tx.SelectTask(taskId); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}
	if task == nil {
		return 0, nil, errors.New("task not found")
	}

	start, end := form.Range()
	file, err := os.OpenFile(task.LogPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, nil, nil
	}
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		if start <= count && count < end {
			line := scanner.Text()
			if line != "" {
				resultList = append(resultList, &results.TaskLogLine{
					LineNum:  count + 1,
					LineText: line,
				})
			}
		}
		count++
	}
	total = count
	return total, resultList, nil
}
