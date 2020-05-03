package constants

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"   // 调度中
	TaskStatusRunning   TaskStatus = "RUNNING"   // 运行中
	TaskStatusFinished  TaskStatus = "FINISHED"  // 已完成
	TaskStatusError     TaskStatus = "ERROR"     // 错误
	TaskStatusCancelled TaskStatus = "CANCELLED" // 取消
)
