package api

type TimeSeriesDataFilter struct {
	ProjectID  *string     `json:"project_id"`
	SensorMac  *string `json:"sensor_mac"`
	SensorType *string `json:"sensor_type"`
	ReceiveNo  *string `json:"receive_no"`
}

type TimeSeriesQueryRequest struct {
	Measurement string               `json:"measurement"`
	Start       string               `json:"start"`
	Stop        string               `json:"stop"`
	Interval    string               `json:"interval"`
	Filter      TimeSeriesDataFilter `json:"filter"`
}
