package task

//
//import (
//	"crawlab-lite/forms"
//	"crawlab-lite/models"
//	"crawlab-lite/services"
//	"errors"
//	"github.com/apex/log"
//	"github.com/robfig/cron/v3"
//	"runtime/debug"
//)
//
//var Scheduler *scheduler
//
//type scheduler struct {
//	cron *cron.Cron
//}
//
//func (s *scheduler) Start() error {
//	exec := cron.New(cron.WithSeconds())
//
//	// 启动cron服务
//	s.cron.Start()
//
//	// 更新调度列表
//	if err := s.Update(); err != nil {
//		log.Errorf("update scheduler error: %s", err.Error())
//		debug.PrintStack()
//		return err
//	}
//
//	// 每30秒更新一次调度列表
//	spec := "*/30 * * * * *"
//	if _, err := exec.AddFunc(spec, updateSchedules); err != nil {
//		log.Errorf("add func update schedulers error: %s", err.Error())
//		debug.PrintStack()
//		return err
//	}
//
//	return nil
//}
//
//func (s *scheduler) Add(schedule models.Schedule) error {
//	spec := schedule.Cron
//
//	// 添加定时调度
//	eid, err := s.cron.AddFunc(spec, scheduleTask(schedule))
//	if err != nil {
//		log.Errorf("add func task error: %s", err.Error())
//		debug.PrintStack()
//		return err
//	}
//
//	// 更新EntryID
//	schedule.EntryId = eid
//
//	// 更新状态
//	schedule.Enabled = true
//
//	// 保存定时调度
//	if err := services.SaveSchedule(&schedule); err != nil {
//		log.Errorf("schedule save error: %s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//// 清除所有定时调度
//func (s *scheduler) RemoveAll() {
//	entries := s.cron.Entries()
//	for i := 0; i < len(entries); i++ {
//		s.cron.Remove(entries[i].ID)
//	}
//}
//
//// 禁用定时调度
//func (s *scheduler) Disable(id string) error {
//	schedule, err := services.QueryScheduleById(id)
//	if err != nil {
//		return err
//	}
//	if schedule.EntryId == 0 {
//		return errors.New("entry id not found")
//	}
//
//	// 从cron服务中删除该调度
//	s.cron.Remove(schedule.EntryId)
//
//	// 更新状态
//	schedule.Enabled = false
//
//	// 保存定时调度
//	if err := services.SaveSchedule(schedule); err != nil {
//		log.Errorf("schedule save error: %s", err.Error())
//		return err
//	}
//
//	return nil
//}
//
//// 启用定时调度
//func (s *scheduler) Enable(id string) error {
//	schedule, err := services.QueryScheduleById(id)
//	if err != nil {
//		return err
//	}
//	if err := s.Add(*schedule); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (s *scheduler) Update() error {
//	// 删除所有定时调度
//	s.RemoveAll()
//
//	// 获取所有定时调度
//	sList, err := services.QueryScheduleList(nil)
//	if err != nil {
//		log.Errorf("get scheduler list error: %s", err.Error())
//		debug.PrintStack()
//		return err
//	}
//
//	user, err := models.GetUserByUsername("admin")
//	if err != nil {
//		log.Errorf("get admin user error: %s", err.Error())
//		return err
//	}
//
//	// 遍历调度列表
//	for i := 0; i < len(sList); i++ {
//		// 单个调度
//		job := sList[i]
//
//		if job.Status == constants.ScheduleStatusStop {
//			continue
//		}
//
//		// 兼容以前版本
//		if job.UserId.Hex() == "" {
//			job.UserId = user.Id
//		}
//
//		// 添加到定时调度
//		if err := s.AddJob(job); err != nil {
//			log.Errorf("add job error: %s, job: %s, cron: %s", err.Error(), job.Name, job.Cron)
//			debug.PrintStack()
//			return err
//		}
//	}
//
//	return nil
//}
//
//// 初始化调度器
//func InitScheduler() error {
//	Scheduler = &scheduler{
//		cron: cron.New(cron.WithSeconds()),
//	}
//	if err := Scheduler.Start(); err != nil {
//		log.Errorf("start scheduler error: %s", err.Error())
//		debug.PrintStack()
//		return err
//	}
//	return nil
//}
//
//func updateSchedules() {
//	if err := Scheduler.Update(); err != nil {
//		log.Errorf(err.Error())
//		return
//	}
//}
//
//func scheduleTask(schedule models.Schedule) func() {
//	return func() {
//		_, _ = services.AddTask(forms.TaskForm{
//			SpiderName:      schedule.SpiderName,
//			SpiderVersionId: schedule.SpiderVersionId,
//			ScheduleId:      schedule.Id,
//			Cmd:             schedule.Cmd,
//		})
//	}
//}
