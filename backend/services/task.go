package services

import (
	"crawlab-lite/constants"
	"crawlab-lite/database"
	"crawlab-lite/forms"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/xujiajun/nutsdb"
	"time"
)

func QueryTaskList(pageNum int, pageSize int) (total int, tasks []*models.Task, err error) {
	start := (pageNum - 1) * pageSize
	end := start + pageSize

	if err := database.KvDB.View(func(tx *nutsdb.Tx) error {
		// 查询区间内的所有任务
		if nodes, err := tx.ZRangeByRank(constants.TaskListBucket, start, end); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var task *models.Task
				if err := json.Unmarshal(node.Value, &task); err != nil {
					return err
				}
				tasks = append(tasks, task)
			}
		}

		// 查询数据总数目
		if total, err = tx.ZCard(constants.TaskListBucket); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, tasks, nil
}

func QueryTaskById(id string) (task *models.Task, err error) {
	if err := database.KvDB.View(func(tx *nutsdb.Tx) error {
		if node, err := tx.ZGetByKey(constants.TaskListBucket, []byte(id)); err != nil {
			if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
				return nil
			}
			return err
		} else if err := json.Unmarshal(node.Value, &task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return task, nil
}

func PopPendingTask() (task *models.Task, err error) {
	if err := database.KvDB.Update(func(tx *nutsdb.Tx) error {
		// 查询爬虫下的所有任务
		if nodes, err := tx.ZMembers(constants.TaskListBucket); err != nil {
			if err == nutsdb.ErrBucket {
				return nil
			}
			return err
		} else {
			for _, node := range nodes {
				var t *models.Task
				if err := json.Unmarshal(node.Value, &t); err != nil {
					return err
				}
				if t.Status == constants.TaskStatusPending {
					t.Status = constants.TaskStatusProcessing
					value, _ := json.Marshal(&t)
					if err := tx.ZAdd(constants.TaskListBucket, []byte(t.Id), float64(node.Score()), value); err != nil {
						return err
					}
					task = t
					return nil
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return task, nil
}

func AddTask(form forms.TaskForm) (task *models.Task, err error) {
	return newTask(form.SpiderName, "", form.Cmd)
}

func newTask(spiderName string, scheduleId string, cmd string) (task *models.Task, err error) {
	// 检查爬虫是否存在
	if spider, err := QuerySpiderByName(spiderName); err != nil {
		return nil, err
	} else if spider == nil {
		return nil, errors.New("spider not found")
	}

	task = &models.Task{
		Id:         uuid.New().String(),
		SpiderName: spiderName,
		ScheduleId: scheduleId,
		Status:     constants.TaskStatusPending,
		Cmd:        cmd,
		CreateTs:   time.Now(),
	}

	// 存储任务信息
	if err := SaveTask(task); err != nil {
		return nil, err
	}

	return task, nil
}

func SaveTask(task *models.Task) (err error) {
	// 存储任务信息
	if err := database.KvDB.Update(func(tx *nutsdb.Tx) error {
		score := utils.ConvertTimestamp(task.CreateTs)
		value, _ := json.Marshal(&task)
		return tx.ZAdd(constants.TaskListBucket, []byte(task.Id), score, value)
	}); err != nil {
		return err
	}

	return nil
}

func UpdateTaskStatus(id string, status constants.TaskStatus) (res interface{}, err error) {
	task, err := QueryTaskById(id)
	if err != nil {
		return nil, err
	}

	task.Status = status

	// 存储任务信息
	if err := SaveTask(task); err != nil {
		return nil, err
	}

	return nil, nil
}
