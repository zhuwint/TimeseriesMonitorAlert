package influxdb

import "time"

// Point use for http response
type Point struct {
	Time     time.Time `json:"time"`
	Value    *float64  `json:"value"`
	FieldTag string    `json:"field_tag"`
}
