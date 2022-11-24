package models

type Task struct {
	TaskId      string `gorm:"column:task_id;primaryKey;not null" json:"task_id"`
	ProjectId   string `gorm:"column:project_id;not null" json:"project_id"`
	TaskType    string `gorm:"column:task_type;not null" json:"task_type"`
	Content     string `gorm:"column:content;not null" json:"content"`
	Created     string `gorm:"column:CREATEDATE;->" json:"created"`
	Updated     string `gorm:"column:UPDATEDATE;->" json:"updated"`
	Description string `gorm:"column:DESCRIPTION;->" json:"description"`
}

func (t Task) TableName() string {
	return "executor_task"
}
