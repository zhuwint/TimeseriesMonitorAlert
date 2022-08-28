package client

import (
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var (
	influxdbClient influxdb2.Client
	once           = sync.Once{}
)

func InitInfluxdbClient(address, bucket, org, token string) {
	once.Do(func() {
		influxdbClient = influxdb2.NewClient(address, token)
	})
}

func Influxdb() influxdb2.Client {
	return influxdbClient
}
