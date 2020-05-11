package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/xujiajun/nutsdb"
	"time"
)

// 查询区间内的所有任务
func (t *Tx) SelectAllTasksLimit(start int, end int) (tasks []*models.Task, err error) {
	if nodes, err := t.tx.ZRangeByRank(constants.TaskListBucket, start, end); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var task *models.Task
			if err = json.Unmarshal(node.Value, &task); err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// 所有任务的总数目
func (t *Tx) CountTasks() (total int, err error) {
	if total, err = t.tx.ZCard(constants.TaskListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}

// 根据 ID 查询任务
func (t *Tx) SelectTask(id uuid.UUID) (task *models.Task, err error) {
	if node, err := t.tx.ZGetByKey(constants.TaskListBucket, []byte(id.String())); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if err = json.Unmarshal(node.Value, &task); err != nil {
		return nil, err
	}
	return task, nil
}

// 根据状态查询任务
func (t *Tx) SelectTaskWhereStatus(status constants.TaskStatus) (task *models.Task, err error) {
	if nodes, err := t.tx.ZMembers(constants.TaskListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var _task *models.Task
			if err = json.Unmarshal(node.Value, &_task); err != nil {
				return nil, err
			}
			if _task.Status == status {
				task = _task
			}
		}
	}
	return task, nil
}

// 插入新任务
func (t *Tx) InsertTask(task *models.Task) (err error) {
	if task.Id == uuid.Nil {
		task.Id = uuid.NewV4()
	}
	if task.CreateTs.IsZero() {
		task.CreateTs = time.Now()
	}
	if task.UpdateTs.IsZero() {
		task.UpdateTs = time.Now()
	}
	if task.Status == "" {
		task.Status = constants.TaskStatusPending
	}

	score := utils.ConvertTimestamp(task.UpdateTs)
	value, _ := json.Marshal(&task)
	if err = t.tx.ZAdd(constants.TaskListBucket, []byte(task.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 更新任务
func (t *Tx) UpdateTask(task *models.Task) (err error) {
	task.UpdateTs = time.Now()
	score := utils.ConvertTimestamp(task.UpdateTs)
	value, _ := json.Marshal(&task)
	if err = t.tx.ZAdd(constants.TaskListBucket, []byte(task.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除任务
func (t *Tx) DeleteTask(id uuid.UUID) (err error) {
	if err = t.tx.ZRem(constants.TaskListBucket, id.String()); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	}
	return nil
}

// 根据爬虫 ID 删除所有任务
func (t *Tx) DeleteAllTasksWhereSpiderId(spiderId uuid.UUID) (err error) {
	if nodes, err := t.tx.ZMembers(constants.TaskListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	} else {
		for _, node := range nodes {
			var task *models.Task
			if err = json.Unmarshal(node.Value, &task); err != nil {
				return err
			}
			if task.SpiderId == spiderId {
				if err = t.tx.ZRem(constants.TaskListBucket, task.Id.String()); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}
