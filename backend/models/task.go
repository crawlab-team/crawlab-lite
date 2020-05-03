package models

import (
	"crawlab-lite/constants"
	"time"
)

type Task struct {
	Id         string               `json:"id"`
	SpiderName string               `json:"spider_name"`
	ScheduleId string               `json:"schedule_id"`
	Status     constants.TaskStatus `json:"status"`
	Cmd        string               `json:"cmd"`
	Error      string               `json:"error"`
	StartTs    time.Time            `json:"start_ts"`
	FinishTs   time.Time            `json:"finish_ts"`
}
