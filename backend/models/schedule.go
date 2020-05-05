package models

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Schedule struct {
	Id              string       `json:"id"`
	SpiderName      string       `json:"spider_name"`
	SpiderVersionId string       `json:"spider_version_id"`
	Cron            string       `json:"cron"`
	EntryId         cron.EntryID `json:"entry_id"`
	Cmd             string       `json:"cmd"`
	Enabled         bool         `json:"enabled"`
	Description     string       `json:"description"`
	CreateTs        time.Time    `json:"create_ts"`
	UpdateTs        time.Time    `json:"update_ts"`
}
