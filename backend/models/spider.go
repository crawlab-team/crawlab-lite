package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Spider struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateTs    time.Time `json:"create_ts"`
	UpdateTs    time.Time `json:"update_ts"`
}

type SpiderVersion struct {
	Id       uuid.UUID `json:"id"`
	MD5      string    `json:"md5"`
	SpiderId uuid.UUID `json:"spider_id"`
	Path     string    `json:"path"`
	CreateTs time.Time `json:"create_ts"`
	UpdateTs time.Time `json:"update_ts"`
}
