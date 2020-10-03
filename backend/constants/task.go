package constants

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"    // 等待调度
	TaskStatusProcessing TaskStatus = "PROCESSING" // 处理中
	TaskStatusRunning    TaskStatus = "RUNNING"    // 运行中
	TaskStatusFinished   TaskStatus = "FINISHED"   // 已完成
	TaskStatusError      TaskStatus = "ERROR"      // 错误
	TaskStatusCancelled  TaskStatus = "CANCELLED"  // 取消
)

const (
	TaskFinish string = "FINISH"
	TaskCancel string = "CANCEL"
)

type TaskLogStd string

const (
	TaskLogStdOut TaskLogStd = "STDOUT"
	TaskLogStdErr TaskLogStd = "STDERR"
)
