package models

type Spider struct {
	Name     string `json:"name"`
	CreateTs int64  `json:"create_ts"`
}

type SpiderVersion struct {
	Id       string `json:"id"`
	Path     string `json:"path"`
	UploadTs int64  `json:"upload_ts"`
}
