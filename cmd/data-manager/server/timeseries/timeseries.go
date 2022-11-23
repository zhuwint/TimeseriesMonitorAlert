package timeseries

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"timeseries/pkg/api"
	influxsvc "timeseries/pkg/service/influxdb"
	"timeseries/pkg/utils/kv"

	"github.com/gin-gonic/gin"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"
const TIME_FORMAT = "2006-01-02T15:04:05Z"

func QueryTimeseries(ctx *gin.Context) {
	var reqBody api.TimeSeriesQueryRequest
	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: err.Error()})
		ctx.Abort()
		return
	}

	var filters []kv.KV
	t := reflect.TypeOf(reqBody.Filter)
	v := reflect.ValueOf(reqBody.Filter)
	for i := 0; i < t.NumField(); i++ {
		value, ok := v.Field(i).Interface().(*string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: "parse body failed"})
			ctx.Abort()
			return
		}
		if value == nil {
			continue
		}
		filters = append(filters, kv.KV{Key: t.Field(i).Tag.Get("json"), Value: *value})
	}

	start, err := time.ParseInLocation(TIME_LAYOUT, reqBody.Start, time.Local)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: "time format error"})
		ctx.Abort()
		return
	}
	stop, err := time.ParseInLocation(TIME_LAYOUT, reqBody.Stop, time.Local)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: "time format error"})
		ctx.Abort()
		return
	}

	query := influxsvc.GeneralQuery{
		Bucket:      influxsvc.BUCKET,
		Measurement: reqBody.Measurement,
		Fields:      []string{"value"},
		Filters:     filters,
		Aggregate: influxsvc.Aggregate{
			Enable: true,
			Every:  reqBody.Interval,
			Fn:     "mean",
		},
		Range: influxsvc.Range{
			Start: start.Format(TIME_FORMAT),
			Stop:  stop.Format(TIME_FORMAT),
		},
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	res, err := influxsvc.Query(query.TransToFlux(), timeoutCtx)
	if err != nil {
		ctx.JSON(http.StatusOK, api.ReplyError{Code: api.InternelError, Msg: err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, api.ReplyJson{Data: res})
}
