package dao

import (
	"crawlab-lite/constants"
	"crawlab-lite/models"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
)

// 查询所有定时调度
func (t *Tx) SelectAllSchedules() (schedules []*models.Schedule, err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return schedules, nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var schedule *models.Schedule
		if err = json.Unmarshal(v, &schedule); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// 根据 ID 查询定时调度
func (t *Tx) SelectSchedule(id uuid.UUID) (schedule *models.Schedule, err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil, nil
	}
	value := b.Get([]byte(id.String()))
	if len(value) == 0 {
		return nil, nil
	}
	if err = json.Unmarshal(value, &schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

// 插入新定时调度
func (t *Tx) InsertSchedule(schedule *models.Schedule) (err error) {
	if schedule.Id == uuid.Nil {
		schedule.Id = uuid.NewV4()
	}
	if schedule.CreateTs.IsZero() {
		schedule.CreateTs = time.Now()
	}
	if schedule.UpdateTs.IsZero() {
		schedule.UpdateTs = time.Now()
	}

	value, err := json.Marshal(&schedule)
	if err != nil {
		return err
	}
	b, err := t.tx.CreateBucketIfNotExists([]byte(constants.ScheduleBucket))
	if err != nil {
		return err
	}
	if err = b.Put([]byte(schedule.Id.String()), value); err != nil {
		return err
	}
	return nil
}

// 更新定时调度
func (t *Tx) UpdateSchedule(schedule *models.Schedule) (err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil
	}
	schedule.UpdateTs = time.Now()
	value, _ := json.Marshal(&schedule)
	if err = b.Put([]byte(schedule.Id.String()), value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除定时调度
func (t *Tx) DeleteSchedule(id uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil
	}
	if err = b.Delete([]byte(id.String())); err != nil {
		return err
	}
	return nil
}

// 根据爬虫 ID 删除所有定时调度
func (t *Tx) DeleteAllSchedulesWhereSpiderId(spiderId uuid.UUID) (err error) {
	b := t.tx.Bucket([]byte(constants.ScheduleBucket))
	if b == nil {
		return nil
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var schedule *models.Schedule
		if err = json.Unmarshal(v, &schedule); err != nil {
			return err
		}
		if schedule.SpiderId == spiderId {
			if err = b.Delete([]byte(schedule.Id.String())); err != nil {
				return err
			}
		}
	}
	return nil
}
