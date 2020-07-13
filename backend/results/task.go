package results

import (
	"crawlab-lite/constants"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Task struct {
	Id              uuid.UUID            `json:"id"`
	SpiderId        uuid.UUID            `json:"spider_id"`
	SpiderName      string               `json:"spider_name"`
	SpiderVersionId uuid.UUID            `json:"spider_version_id"`
	ScheduleId      uuid.UUID            `json:"schedule_id"`
	Status          constants.TaskStatus `json:"status"`
	Cmd             string               `json:"cmd"`
	Error           string               `json:"error"`
	CreateTs        time.Time            `json:"create_ts"`
	UpdateTs        time.Time            `json:"update_ts"`
	StartTs         time.Time            `json:"start_ts"`
	FinishTs        time.Time            `json:"finish_ts"`
	ResultCount     int                  `json:"result_count"`
	ErrorLogCount   int                  `json:"error_log_count"`
	WaitDuration    float64              `json:"wait_duration"`
	RuntimeDuration float64              `json:"runtime_duration"`
	TotalDuration   float64              `json:"total_duration"`
}

type TaskLogLine struct {
	LineNum  int    `json:"line_num"`
	LineText string `json:"line_text"`
}
