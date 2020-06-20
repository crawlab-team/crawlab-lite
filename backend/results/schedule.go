package results

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Schedule struct {
	Id              uuid.UUID `json:"id"`
	SpiderId        uuid.UUID `json:"spider_id"`
	SpiderName      string    `json:"spider_name"`
	SpiderVersionId uuid.UUID `json:"spider_version_id"`
	Cron            string    `json:"cron"`
	Cmd             string    `json:"cmd"`
	Enabled         bool      `json:"enabled"`
	Description     string    `json:"description"`
	CreateTs        time.Time `json:"create_ts"`
	UpdateTs        time.Time `json:"update_ts"`
}
