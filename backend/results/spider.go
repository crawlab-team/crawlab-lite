package results

import (
	"crawlab-lite/constants"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Spider struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	CreateTs    time.Time            `json:"create_ts"`
	UpdateTs    time.Time            `json:"update_ts"`
	LastRunTs   time.Time            `json:"last_run_ts"` // 最后一次执行时间
	LastStatus  constants.TaskStatus `json:"last_status"` // 最后执行状态
}

type SpiderVersion struct {
	Id       uuid.UUID `json:"id"`
	FileHash string    `json:"file_hash"`
	SpiderId uuid.UUID `json:"spider_id"`
	Path     string    `json:"path"`
	CreateTs time.Time `json:"create_ts"`
	UpdateTs time.Time `json:"update_ts"`
}
