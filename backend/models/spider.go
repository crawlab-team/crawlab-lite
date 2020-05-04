package models

import (
	"time"
)

type Spider struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	CreateTs time.Time `json:"create_ts"`
	UpdateTs time.Time `json:"update_ts"`
}

type SpiderVersion struct {
	Id         string    `json:"id"`
	SpiderName string    `json:"spider_name"`
	Path       string    `json:"path"`
	CreateTs   time.Time `json:"create_ts"`
	UpdateTs   time.Time `json:"update_ts"`
}
