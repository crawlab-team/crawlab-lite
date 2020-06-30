package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
)

// 查询所有任务
func (t *Tx) SelectAllTasks() (tasks []*models.Task, err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return tasks, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var task *models.Task
		if err = json.Unmarshal(v, &task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// 根据爬虫 ID 查询任务
func (t *Tx) SelectTasksWhereSpiderId(spiderId uuid.UUID) (tasks []*models.Task, err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return tasks, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var task *models.Task
		if err = json.Unmarshal(v, &task); err != nil {
			return nil, err
		}
		if task.SpiderId == spiderId {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// 根据 ID 查询任务
func (t *Tx) SelectTask(id uuid.UUID) (task *models.Task, err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(id.String()))
	if len(value) == 0 {
		return nil, nil
	}
	if err = json.Unmarshal(value, &task); err != nil {
		return nil, err
	}
	return task, nil
}

// 根据状态查询任务
func (t *Tx) SelectTaskWhereStatus(status constants.TaskStatus) (task *models.Task, err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return nil, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if err = json.Unmarshal(v, &task); err != nil {
			return nil, err
		}
		if task.Status == status {
			return task, nil
		}
	}
	return nil, nil
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

	value, err := json.Marshal(&task)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.TaskBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(task.Id.String()), value); err != nil {
		return err
	}
	return nil
}

// 更新任务
func (t *Tx) UpdateTask(task *models.Task) (err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return nil
	}
	task.UpdateTs = time.Now()
	value, _ := json.Marshal(&task)
	if err = b.Put([]byte(task.Id.String()), value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除任务
func (t *Tx) DeleteTask(id uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(id.String())); err != nil {
		return err
	}
	return nil
}

// 根据爬虫 ID 删除所有任务
func (t *Tx) DeleteTasksWhereSpiderId(spiderId uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.TaskBucket))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var task *models.Task
		if err = json.Unmarshal(v, &task); err != nil {
			return err
		}
		if task.SpiderId == spiderId {
			if err = b.Delete([]byte(task.Id.String())); err != nil {
				return err
			}
		}
	}
	return nil
}
