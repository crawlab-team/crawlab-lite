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

// 查询区间内的所有定时调度
func (t *Tx) SelectAllSchedulesLimit(start int, end int) (schedules []*models.Schedule, err error) {
	if nodes, err := t.tx.ZRangeByRank(constants.ScheduleListBucket, -(start + 1), -(end + 1)); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var schedule *models.Schedule
			if err = json.Unmarshal(node.Value, &schedule); err != nil {
				return nil, err
			}
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// 查询所有定时调度
func (t *Tx) SelectAllSchedules() (schedules []*models.Schedule, err error) {
	if nodes, err := t.tx.ZMembers(constants.ScheduleListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil, nil
		}
		return nil, err
	} else {
		for _, node := range nodes {
			var schedule *models.Schedule
			if err = json.Unmarshal(node.Value, &schedule); err != nil {
				return nil, err
			}
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// 所有定时调度的总数目
func (t *Tx) CountSchedules() (total int, err error) {
	if total, err = t.tx.ZCard(constants.ScheduleListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}

// 根据 ID 查询定时调度
func (t *Tx) SelectSchedule(id uuid.UUID) (schedule *models.Schedule, err error) {
	if node, err := t.tx.ZGetByKey(constants.ScheduleListBucket, []byte(id.String())); err != nil {
		if err == nutsdb.ErrBucket || err == nutsdb.ErrNotFoundKey {
			return nil, nil
		}
		return nil, err
	} else if err = json.Unmarshal(node.Value, &schedule); err != nil {
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

	score := utils.ConvertTimestamp(schedule.UpdateTs)
	value, _ := json.Marshal(&schedule)
	if err = t.tx.ZAdd(constants.ScheduleListBucket, []byte(schedule.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 更新定时调度
func (t *Tx) UpdateSchedule(schedule *models.Schedule) (err error) {
	schedule.UpdateTs = time.Now()
	score := utils.ConvertTimestamp(schedule.UpdateTs)
	value, _ := json.Marshal(&schedule)
	if err = t.tx.ZAdd(constants.ScheduleListBucket, []byte(schedule.Id.String()), score, value); err != nil {
		return err
	}
	return nil
}

// 通过 ID 删除定时调度
func (t *Tx) DeleteSchedule(id uuid.UUID) (err error) {
	if err = t.tx.ZRem(constants.ScheduleListBucket, id.String()); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	}
	return nil
}

// 根据爬虫 ID 删除所有定时调度
func (t *Tx) DeleteAllSchedulesWhereSpiderId(spiderId uuid.UUID) (err error) {
	if nodes, err := t.tx.ZMembers(constants.ScheduleListBucket); err != nil {
		if err == nutsdb.ErrBucket {
			return nil
		}
		return err
	} else {
		for _, node := range nodes {
			var schedule *models.Schedule
			if err = json.Unmarshal(node.Value, &schedule); err != nil {
				return err
			}
			if schedule.SpiderId == spiderId {
				if err = t.tx.ZRem(constants.ScheduleListBucket, schedule.Id.String()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
