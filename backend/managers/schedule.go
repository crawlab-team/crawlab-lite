package managers

import (
	"crawlab-lite/dao"
	"crawlab-lite/database"
	"crawlab-lite/models"
	"github.com/apex/log"
	"github.com/robfig/cron/v3"
)

// 任务调度器
var Scheduler *scheduler

type scheduler struct {
	cron *cron.Cron
}

func (s *scheduler) Start() error {
	// 启动cron服务
	s.cron.Start()

	// 刷新调度列表
	if err := s.Flush(); err != nil {
		log.Errorf("update scheduler error: %s", err.Error())
		return err
	}

	return nil
}

func (s *scheduler) Add(schedule *models.Schedule) (err error) {
	s.Remove(schedule)

	// 添加定时调度
	spec := schedule.Cron
	if schedule.EntryId, err = s.cron.AddFunc(spec, scheduleTask(*schedule)); err != nil {
		log.Errorf("add func task error: %s", err.Error())
		return err
	}

	return nil
}

func (s *scheduler) Remove(schedule *models.Schedule) {
	if schedule.EntryId != 0 {
		s.cron.Remove(schedule.EntryId)
		schedule.EntryId = 0
	}
}

// 刷新所有调度器
func (s *scheduler) Flush() error {
	// 删除所有定时调度
	entries := s.cron.Entries()
	for i := 0; i < len(entries); i++ {
		s.cron.Remove(entries[i].ID)
	}

	if err := dao.WriteTx(database.MainDB, func(tx dao.Tx) error {
		if schedules, err := tx.SelectAllSchedules(); err != nil {
			return err
		} else {
			for _, schedule := range schedules {
				if schedule.Cron == "" || schedule.Enabled == false {
					continue
				}
				// 添加到定时调度
				if err = s.Add(schedule); err != nil {
					log.Errorf("add scheduler error: %s, scheduler: %s, cron: %s", err.Error(), schedule.Id, schedule.Cron)
					return err
				} else {
					if err = tx.UpdateSchedule(schedule); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}); err != nil {
		log.Errorf("flush schedulers error: %s", err.Error())
		return err
	}

	return nil
}

// 初始化调度器
func InitScheduler() error {
	Scheduler = &scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
	if err := Scheduler.Start(); err != nil {
		return err
	}
	return nil
}

func scheduleTask(schedule models.Schedule) func() {
	return func() {
		task := &models.Task{
			SpiderId:        schedule.SpiderId,
			SpiderVersionId: schedule.SpiderVersionId,
			ScheduleId:      schedule.Id,
			Cmd:             schedule.Cmd,
		}
		if err := dao.WriteTx(database.MainDB, func(tx dao.Tx) error {
			if err := tx.InsertTask(task); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return
		}
	}
}
