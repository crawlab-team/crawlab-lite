package models

import "time"

type Spider struct {
	Name     string    `json:"name"`
	CreateTs time.Time `json:"create_ts"`
}

type SpiderVersion struct {
	Id       string    `json:"id"`
	Path     string    `json:"path"`
	CreateTs time.Time `json:"create_ts"`
}
