package services

import (
	"crawlab-lite/constants"
	"crawlab-lite/dao"
	"crawlab-lite/forms"
	"crawlab-lite/managers"
	"crawlab-lite/models"
	"crawlab-lite/results"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
)

func QuerySchedulePage(page forms.PageForm) (total int, resultList []*results.Schedule, err error) {
	start, end := page.Range()

	if err := dao.ReadTx(func(tx dao.Tx) error {
		schedules, err := tx.SelectAllSchedulesLimit(start, end)
		if err != nil {
			return err
		}
		cache := map[uuid.UUID]*models.Spider{}
		for _, schedule := range schedules {
			var result results.Schedule
			if err := copier.Copy(&result, schedule); err != nil {
				return err
			}
			var spider *models.Spider
			spider, exists := cache[schedule.SpiderId]
			if !exists {
				if spider, err = tx.SelectSpider(schedule.SpiderId); err != nil {
					return err
				}
			}
			if spider == nil {
				return errors.New("spider not found")
			}
			result.SpiderName = spider.Name
			resultList = append(resultList, &result)
		}
		if total, err = tx.CountSchedules(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, resultList, nil
}

func QueryScheduleById(id uuid.UUID) (result *results.Schedule, err error) {
	if err := dao.ReadTx(func(tx dao.Tx) error {
		schedule, err := tx.SelectSchedule(id)
		if err != nil {
			return err
		}
		if err := copier.Copy(&result, schedule); err != nil {
			return err
		}
		spider, err := tx.SelectSpider(schedule.SpiderId)
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

func AddSchedule(form forms.ScheduleCreateForm) (result *results.Schedule, err error) {
	if form.Cron != "" && CheckCron(form.Cron) == false {
		return nil, errors.New("schedule cron is invalid")
	}
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

		schedule := &models.Schedule{
			SpiderId:        form.SpiderId,
			SpiderVersionId: form.SpiderVersionId,
			Cron:            form.Cron,
			Cmd:             form.Cmd,
			Enabled:         true,
		}
		// 添加定时
		if err = managers.Scheduler.Add(schedule); err != nil {
			return err
		}
		// 存储任务信息
		if err := tx.InsertSchedule(schedule); err != nil {
			return err
		}
		if err := copier.Copy(&result, schedule); err != nil {
			return err
		}
		spider, err := tx.SelectSpider(schedule.SpiderId)
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

func ModifySchedule(id uuid.UUID, form forms.ScheduleUpdateForm) (result *results.Schedule, err error) {
	if form.Cron != "" && CheckCron(form.Cron) == false {
		return nil, errors.New("schedule cron is invalid")
	}
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查调度是否存在
		schedule, err := tx.SelectSchedule(id)
		if err != nil {
			return err
		}
		if schedule == nil {
			return errors.New("schedule not found")
		}
		isEnabled := form.Enabled == constants.Enabled
		// 如果更新调度的可用状态，同步增删定时
		if form.Enabled != 0 {
			if isEnabled {
				if err = managers.Scheduler.Add(schedule); err != nil {
					return err
				}
			} else {
				managers.Scheduler.Remove(schedule)
			}
		}
		if form.Cmd != "" {
			schedule.Cmd = form.Cmd
		}
		if form.Cron != "" {
			schedule.Cron = form.Cron
		}
		if form.Enabled != 0 {
			schedule.Enabled = isEnabled
		}
		if form.Description != "" {
			schedule.Description = form.Description
		}
		// 更新调度信息
		if err := tx.UpdateSchedule(schedule); err != nil {
			return err
		}
		if err := copier.Copy(&result, schedule); err != nil {
			return err
		}
		spider, err := tx.SelectSpider(schedule.SpiderId)
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

func RemoveSchedule(id uuid.UUID) (res interface{}, err error) {
	if err := dao.WriteTx(func(tx dao.Tx) error {
		// 检查调度是否存在
		if schedule, err := tx.SelectSchedule(id); err != nil {
			return err
		} else if schedule == nil {
			return errors.New("schedule not found")
		} else {
			// 删除调度
			if err = tx.DeleteSchedule(id); err != nil {
				return err
			}
			// 清除定时
			if schedule.EntryId != 0 {
				managers.Scheduler.Remove(schedule)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

// 检查 cron 表达式是否正确
func CheckCron(spec string) bool {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	if _, err := parser.Parse(spec); err != nil {
		return false
	}
	return true
}
