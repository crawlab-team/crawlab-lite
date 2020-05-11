package models

import (
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Schedule struct {
	Id              uuid.UUID    `json:"id"`
	SpiderId        uuid.UUID    `json:"spider_id"`
	SpiderVersionId uuid.UUID    `json:"spider_version_id"`
	Cron            string       `json:"cron"`
	EntryId         cron.EntryID `json:"entry_id"`
	Cmd             string       `json:"cmd"`
	Enabled         bool         `json:"enabled"`
	Description     string       `json:"description"`
	CreateTs        time.Time    `json:"create_ts"`
	UpdateTs        time.Time    `json:"update_ts"`
}
