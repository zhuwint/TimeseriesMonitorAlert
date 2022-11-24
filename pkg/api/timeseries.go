package api

import (
	"fmt"
)

type TimeSeriesDataFilter struct {
	ProjectID  *string `json:"project_id"`
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

// UnvariedSeries 单变量时间序列查询
type UnvariedSeries struct {
	TimeSeriesDataFilter `json:",inline"`
}

func (u UnvariedSeries) Validate() error {
	if u.ProjectID == nil || u.SensorMac == nil || u.ReceiveNo == nil || u.SensorType == nil {
		return fmt.Errorf("information cannot be empty")
	}
	if *u.ProjectID == "" || *u.SensorMac == "" || *u.ReceiveNo == "" || *u.SensorType == "" {
		return fmt.Errorf("information cannot be empty")
	}
	return nil
}
