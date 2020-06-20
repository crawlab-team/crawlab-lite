package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"crawlab-lite/utils"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/xujiajun/nutsdb"
	"math"
	"time"
)

// 查询区间内的所有任务
func (t *Tx) SelectAllTasksLimit(start int, end int) (tasks []*models.Task, err error) {
	if nodes, err := t.tx.ZRangeByRank(constants.TaskListBucket, -(start + 1), -(end + 1)); err != nil {
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

// 根据爬虫 ID 查询所有任务
func (t *Tx) SelectTasksWhereSpiderIdLimit(spiderId uuid.UUID, start int, end int) (tasks []*models.Task, err error) {
	if total, err := t.CountTasks(); err != nil {
		return nil, err
	} else {
		if nodes, err := t.tx.ZRangeByRank(constants.TaskListBucket, -1, -total); err != nil {
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
				if task.SpiderId == spiderId {
					tasks = append(tasks, task)
				}
			}
			if len(tasks) < start {
				return nil, nil
			}
			end = int(math.Min(float64(end), float64(len(tasks))))
			return tasks[start:end], nil
		}
	}
}

// 根据爬虫 ID 查询所有任务的总数目
func (t *Tx) CountTasks() (total int, err error) {
	if total, err = t.tx.ZCard(constants.TaskListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}

// 所有任务的总数目
func (t *Tx) CountTasksBySpiderId(spiderId uuid.UUID) (total int, err error) {
	if total, err := t.CountTasks(); err != nil {
		return 0, err
	} else {
		if nodes, err := t.tx.ZRangeByRank(constants.TaskListBucket, -1, -total); err != nil {
			if err == nutsdb.ErrBucket {
				return 0, nil
			}
			return 0, err
		} else {
			tasks := make([]*models.Task, 0, len(nodes))
			for _, node := range nodes {
				var task *models.Task
				if err = json.Unmarshal(node.Value, &task); err != nil {
					return 0, err
				}
				if task.SpiderId == spiderId {
					tasks = append(tasks, task)
				}
			}
			return len(tasks), nil
		}
	}
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
				return _task, nil
			}
		}
	}
	return task, nil
}

// 根据爬虫 ID 查询最近的任务
func (t *Tx) SelectLatestTaskWhereSpiderId(spiderId uuid.UUID) (task *models.Task, err error) {
	if total, err := t.CountTasks(); err != nil {
		return nil, err
	} else {
		if nodes, err := t.tx.ZRangeByRank(constants.TaskListBucket, -1, -total); err != nil {
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
				if _task.SpiderId == spiderId {
					return _task, nil
				}
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
func (t *Tx) DeleteTasksWhereSpiderId(spiderId uuid.UUID) (err error) {
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
			}
		}
	}
	return nil
}
