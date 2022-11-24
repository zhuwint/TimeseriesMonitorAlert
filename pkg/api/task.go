package api

type TaskType string

const (
	TaskTypeStream TaskType = "stream"
	ETaskTypeBatch TaskType = "batch"
)

type Task interface {
	TaskId() string
	TaskType() TaskType
}

type TaskInfo struct {
	Id   string   `json:"task_id"`
	Type TaskType `json:"task_type"`
}

func (t TaskInfo) TaskId() string {
	return t.Id
}

func (t TaskInfo) TaskType() TaskType {
	return t.Type
}

type StreamTaskInfo struct {
	TaskInfo `json:",inline"`
}

type BatchTaskInfo struct {
	TaskInfo    `json:",inline"`
	Target      UnvariedSeries   `json:"target"`      // 目标检测序列
	Independent []UnvariedSeries `json:"independent"` // 其它序列（自变量）
	// DetectModel
}
