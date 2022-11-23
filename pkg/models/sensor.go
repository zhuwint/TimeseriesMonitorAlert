package models

// SensorLocation 只读
type SensorLocation struct {
	ID          int    `gorm:"column:ID;primaryKey;not null;->" json:"id"`
	SensorMac   string `gorm:"column:SENSOR_MAC;->" json:"sensor_mac"`
	Status      int    `gorm:"column:STATUS;not null;->" json:"status"`
	ProjectId   int    `gorm:"column:PROJECT_ID;->" json:"project_id"`
	TypeId      int    `gorm:"column:TYPE_ID;->" json:"type_id"`
	Location1Id int    `gorm:"column:LOCATION_1_ID;->" json:"location_1_id"`
	Location2Id int    `gorm:"column:LOCATION_2_ID;->" json:"location_2_id"`
	Location3Id int    `gorm:"column:LOCATION_3_ID;->" json:"location_3_id"`
	Location4Id int    `gorm:"column:LOCATION_4_ID;->" json:"location_4_id"`
	Created     string `gorm:"column:CREATEDATE;->" json:"created"`
	Updated     string `gorm:"column:UPDATEDATE;->" json:"update"`
	Description string `gorm:"column:DESCRIPTION;->" json:"description"`
	UserId      string `gorm:"column:USERID;->" json:"user_id"`
}

func (s SensorLocation) TableName() string {
	return "sensor_location"
}

// SensorType 只读
type SensorType struct {
	ID          int    `gorm:"column:ID;primaryKey;not null;->" json:"id"`
	TypeName    string `gorm:"column:TYPE_NAME;->" json:"type_name"`
	Created     string `gorm:"column:CREATEDATE;->" json:"created"`
	Updated     string `gorm:"column:UPDATEDATE;->" json:"update"`
	Description string `gorm:"column:DESCRIPTION;->" json:"description"`
}

func (s SensorType) TableName() string {
	return "sensor_type"
}

// SensorGatherType 只读
type SensorGatherType struct {
	ID             int    `gorm:"column:id;primaryKey;not null;->" json:"id"`
	SensorTypeId   int    `gorm:"column:sensor_type_id;not null;->" json:"sensor_type_id"`
	GatherType     string `gorm:"column:gather_type;not null;->" json:"gather_type"`
	ReceiveNumber  int    `gorm:"column:receive_no;not null;->" json:"receive_no"`
	GatherTypeName string `gorm:"column:gather_type_name" json:"gather_type_name"`
	Unit           string `gorm:"column:unit;->" json:"unit"`
	Figures        int    `gorm:"figures;->" json:"figures"`
	Description    string `gorm:"column:DESCRIPTION;->" json:"description"`
}

func (s SensorGatherType) TableName() string {
	return "sensor_gather_type"
}

type ProjectIdName struct {
	ProjectId   int    `gorm:"column:PROJECT_ID;primaryKey;not null;->" json:"project_id"`
	ProjectName string `gorm:"column:PROJECT_NAME;->" json:"project_name"`
	Created     string `gorm:"column:CREATEDATE;->" json:"created"`
	Updated     string `gorm:"column:UPDATEDATE;->" json:"update"`
}

func (p ProjectIdName) TableName() string {
	return "project_id_name"
}

type SiteLocationName struct {
	ID                  int     `gorm:"column:ID;->;primaryKey;not null" json:"id"`
	ProjectId           int     `gorm:"column:project_id;->" json:"project_id"`
	Location1Id         int     `gorm:"column:LOCATION_1_ID;->" json:"location_1_id"`
	Location2Id         int     `gorm:"column:LOCATION_2_ID;->" json:"location_2_id"`
	Location3Id         int     `gorm:"column:LOCATION_3_ID;->" json:"location_3_id"`
	Location4Id         int     `gorm:"column:LOCATION_4_ID;->" json:"location_4_id"`
	LocationName        string  `gorm:"column:LOCATION_NAME;->" json:"location_name"`
	LocationDescription string  `gorm:"column:LOCATION_DESCRIPTION;->" json:"location_description"`
	LocationX           float32 `gorm:"column:LOCATION_X;->" json:"location_x"`
	LocationY           float32 `gorm:"column:LOCATION_Y;->" json:"location_y"`
	Created             string  `gorm:"column:CREATEDATE;->" json:"created"`
	Updated             string  `gorm:"column:UPDATEDATE;->" json:"update"`
	UserId              string  `gorm:"column:USERID;->" json:"user_id"`
}

func (s SiteLocationName) TableName() string {
	return "site_location_name"
}
