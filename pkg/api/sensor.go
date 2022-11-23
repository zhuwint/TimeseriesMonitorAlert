package api

type ProjectResp struct {
	ProjectId   int    `json:"project_id"`
	ProjectName string `json:"project_name"`
}

type SensorReq struct {
	ProjectId *string `json:"project_id" filter:"project_id"`
	Location1 *string `json:"location1" filter:"location_1_id"`
	Location2 *string `json:"location2" filter:"location_2_id"`
	Location3 *string `json:"location3" filter:"location_3_id"`
	Location4 *string `json:"location4" filter:"location_4_id"`
}

type SensorResp struct {
	SensorMac string `json:"sensor_mac"`
	TypeId    int    `json:"type_id"`
	TypeName  string `json:"type_name"`
	Location1 string `json:"location1"`
	Location2 string `json:"location2"`
	Location3 string `json:"location3"`
	Location4 string `json:"location4"`
}

type LocationNode struct {
	Level        int             `json:"level"`
	Location1Id  int             `json:"location_1_id"`
	Location2Id  int             `json:"location_2_id"`
	Location3Id  int             `json:"location_3_id"`
	Location4Id  int             `json:"location_4_id"`
	LocationName string          `json:"location_name"`
	Children     []*LocationNode `json:"children"`
}

type MeasurementResp struct {
	ReceiveNo      int    `json:"receive_no"`
	GatherType     string `json:"gather_type"`
	GatherTypeName string `json:"gather_type_name"`
	Unit           string `json:"unit"`
}
