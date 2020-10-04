package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"encoding/binary"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"go.etcd.io/bbolt"
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
func (t *Tx) SelectFirstTaskWhereStatus(status constants.TaskStatus) (task *models.Task, err error) {
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

// 根据任务 ID 查询所有日志
func (t *Tx) SelectAllTaskLogs(taskId uuid.UUID) (logs []*models.TaskLog, err error) {
	b := t.tx.Bucket([]byte(constants.TaskLogBucket + taskId.String()))
	if b == nil {
		return nil, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var taskLog models.TaskLog
		if err = json.Unmarshal(v, &taskLog); err != nil {
			return nil, err
		}
		logs = append(logs, &taskLog)
	}
	return logs, nil
}

// 根据任务 ID 范围查询日志
func (t *Tx) SelectTaskLogsLimit(taskId uuid.UUID, limit int, offset int) (logs []*models.TaskLog, err error) {
	b := t.tx.Bucket([]byte(constants.TaskLogBucket + taskId.String()))
	if b == nil {
		return nil, nil
	}
	c := b.Cursor()
	count := 0
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if limit > len(logs) && count >= offset {
			var taskLog models.TaskLog
			if err = json.Unmarshal(v, &taskLog); err != nil {
				return nil, err
			}
			logs = append(logs, &taskLog)
		}
		count++
	}
	return logs, nil
}

// 根据任务 ID 查询日志总数
func (t *Tx) CountTaskLogs(taskId uuid.UUID) (count int, err error) {
	b := t.tx.Bucket([]byte(constants.TaskLogBucket + taskId.String()))
	if b == nil {
		return 0, nil
	}
	c := b.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		count++
	}
	return count, nil
}

// 插入新任务日志
func (t *Tx) InsertTaskLog(taskLog *models.TaskLog) (err error) {
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.TaskLogBucket + taskLog.TaskId.String()))
	if err != nil {
		return err
	}
	if taskLog.CreateTs.IsZero() {
		taskLog.CreateTs = time.Now()
	}
	value, err := json.Marshal(&taskLog)
	if err != nil {
		return err
	}
	id, _ := b.NextSequence()
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, id)
	if err = b.Put(bytes, value); err != nil {
		return err
	}
	return nil
}

// 根据任务 ID 删除所有日志
func (t *Tx) DeleteAllTaskLogs(taskId uuid.UUID) (err error) {
	bucket := []byte(constants.TaskLogBucket + taskId.String())
	b := t.tx.Bucket(bucket)
	if b == nil {
		return nil
	}
	if err = b.ForEach(func(k, v []byte) error {
		return b.Delete(k)
	}); err != nil {
		return err
	}
	if err = t.tx.DeleteBucket(bucket); err != nil {
		if err == bbolt.ErrBucketNotFound {
			return nil
		}
		return err
	}
	return nil
}

// 根据任务 ID 删除超过 N 天的日志
func (t *Tx) DeleteTaskLogsOlderThan(taskId uuid.UUID, days int) (err error) {
	expireDate := time.Now().AddDate(0, 0, -days)
	b := t.tx.Bucket([]byte(constants.TaskLogBucket + taskId.String()))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var taskLog models.TaskLog
		if err = json.Unmarshal(v, &taskLog); err != nil {
			return err
		}
		if taskLog.CreateTs.Before(expireDate) {
			if err = b.Delete(k); err != nil {
				return err
			}
		}
	}
	return nil
}
