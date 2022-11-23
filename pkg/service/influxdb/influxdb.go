package influxdb

import (
	"context"
	"fmt"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"golang.org/x/sync/semaphore"
)

type Account struct {
	Address string `yaml:"address"`
	Bucket  string `yaml:"bucket"`
	Token   string `yaml:"token"`
	Org     string `yaml:"org"`
}

func (i Account) Validate() error {
	if i.Bucket == "" {
		return fmt.Errorf("influxdb: bucket coult not be empty")
	}
	if i.Token == "" {
		return fmt.Errorf("influxdb: token could not be empty")
	}
	if i.Org == "" {
		return fmt.Errorf("influxdb: org could not be empty")
	}
	if i.Address == "" {
		return fmt.Errorf("influxdb: address could not be empty")
	}
	return nil
}

type Connector struct {
	Address string
	Bucket  string
	Token   string
	Org     string
	client  influxdb2.Client
	sema    *semaphore.Weighted // use for concurrency
}

var (
	influxClient *Connector
	influxOnce   = &sync.Once{}
)

var (
	BUCKET string
)

func newConnector(address, bucket, token, org string) (*Connector, error) {
	return &Connector{
		Address: address,
		Bucket:  bucket,
		Token:   token,
		Org:     org,
		client:  influxdb2.NewClient(address, token),
		sema:    semaphore.NewWeighted(10),
	}, nil
}

func InitInfluxClient(c Account, bucket string) {
	influxOnce.Do(func() {
		conn, err := newConnector(c.Address, c.Bucket, c.Token, c.Org)
		if err != nil {
			panic(fmt.Errorf("init influx conn failed %s", err.Error()))
		}
		influxClient = conn
	})
	BUCKET = bucket
}

func GetClient() influxdb2.Client {
	if influxClient == nil {
		panic("influx client not init yeat")
	}
	return influxClient.client
}

func CloseClient() {
	GetClient().Close()
}

func Query(script string, ctx context.Context) ([]*Point, error) {
	queryApi := GetClient().QueryAPI(influxClient.Org)
	raw, err := queryApi.Query(ctx, script)
	if err != nil {
		return nil, err
	}

	// must init with size, otherwise gin will return null for empty array
	result := make([]*Point, 0)

	for raw.Next() {
		var p = &Point{
			Time:     raw.Record().Time(),
			Value:    nil,
			FieldTag: raw.Record().ValueByKey("sensor_type").(string),
		}

		value := raw.Record().ValueByKey("_value")
		switch v := value.(type) {
		case float32:
			_v := float64(v)
			p.Value = &_v
		case float64:
			p.Value = &v
		default:
			if value != nil {
				return nil, fmt.Errorf("invalid value type")
			}
		}
		result = append(result, p)
	}

	// sort.Slice(result, func(i, j int) bool {
	// 	return result[i].Time.Before(result[j].Time)
	// })

	if raw.Err() != nil {
		return nil, fmt.Errorf("query parsing error: %s", raw.Err().Error())
	}

	return result, nil
}
