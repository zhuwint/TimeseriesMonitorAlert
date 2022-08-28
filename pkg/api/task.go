package api

type ExecuteTaskType string

const (
	ExecuteTaskTypeStream ExecuteTaskType = "stream"
	ExecuteTaskTypeBatch  ExecuteTaskType = "batch"
)

type ExecuteTask interface {
	Id() string
	Type() ExecuteTaskType
}
