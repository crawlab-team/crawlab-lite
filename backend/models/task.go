package models

import (
	"crawlab-lite/constants"
	"time"
)

type Task struct {
	Id              string               `json:"id"`
	SpiderName      string               `json:"spider_name"`
	SpiderVersionId string               `json:"spider_version_id"`
	ScheduleId      string               `json:"schedule_id"`
	Status          constants.TaskStatus `json:"status"`
	Cmd             string               `json:"cmd"`
	Error           string               `json:"error"`
	CreateTs        time.Time            `json:"create_ts"`
	StartTs         time.Time            `json:"start_ts"`
	FinishTs        time.Time            `json:"finish_ts"`
	ResultCount     int                  `json:"result_count"`
	ErrorLogCount   int                  `json:"error_log_count"`
	WaitDuration    float64              `json:"wait_duration"`
	RuntimeDuration float64              `json:"runtime_duration"`
	TotalDuration   float64              `json:"total_duration"`
}
