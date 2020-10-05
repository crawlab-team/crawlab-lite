package results

import uuid "github.com/satori/go.uuid"

type ResponseBody struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type BatchCount struct {
	SuccessCount int
	FailCount    int
	FailReasons  []map[uuid.UUID]string
}
