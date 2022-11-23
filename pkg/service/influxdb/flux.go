package influxdb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"timeseries/pkg/utils/kv"
)

// InfluxQuery : influxdb query interface
type InfluxQuery interface {
	TransToFlux() string // transform struct to flux script
}

type Filter interface {
}

// BaseQuery : use for unvaried time series query
// use []BaseQuery to query multivariable time series
// filters: string, like
//
//	["host=ubuntu", "cpu=cpu0,cpu=cpu1"]
//
// must provide at least one field
type BaseQuery struct {
	Alias       string   `json:"alias"` // name of the query, use for table join. must be unique
	Bucket      string   `json:"bucket"`
	Measurement string   `json:"measurement"`
	Fields      []string `json:"fields"`
	Filters     []kv.KV  `json:"filters"`
}

func (q BaseQuery) Validate() error {
	if q.Alias == "" {
		return fmt.Errorf("alias could not be empty")
	}
	if q.Bucket == "" {
		return fmt.Errorf("bucket could not be empty")
	}
	if q.Measurement == "" {
		return fmt.Errorf("measurement could not be empty")
	}
	if len(q.Fields) == 0 {
		return fmt.Errorf("must provide at least one field")
	}
	for _, f := range q.Filters {
		if f.Key == "" || f.Value == "" {
			return fmt.Errorf("filter key(value) could not be empty")
		}
	}
	return nil
}

// Aggregate : aggregate with interval and aggregate function
type Aggregate struct {
	Enable      bool   `json:"enable"`
	Every       string `json:"every"`
	Fn          string `json:"fn"`
	CreateEmpty bool   `json:"create_empty"`
}

func (a Aggregate) Validate() error {
	if !a.Enable {
		return nil
	}
	if a.Fn != "mean" && a.Fn != "median" {
		return fmt.Errorf("fn must in [mean, median]")
	}
	if d, err := time.ParseDuration(a.Every); err != nil {
		return err
	} else {
		now := time.Now()
		if !now.Add(d).After(now) {
			return fmt.Errorf("every must be positive time duration")
		}
	}
	return nil
}

// Range : query with time range
type Range struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

func (r Range) Validate() error {
	if r.Start == "" || r.Stop == "" {
		return fmt.Errorf("range parse failed: could not be empty string")
	}

	// start can be all supported unix duration unit relative to stop or absolute time.
	// for example -20h5m3s, 2019-08-28T22:00:00Z.
	// stop can be all supported unix duration unit relative to now or absolute time or now().
	// for example -20h5m3s, 2019-08-28T22:00:00Z, now().
	// supported duration-types are "ns", "us" (or "µs"), "ms", "s", "m", "h".

	start, err := CheckTimeBeforeNow(r.Start)
	if err != nil {
		return fmt.Errorf("range (start) parse failed: %s", err.Error())
	}

	stop, err := CheckTimeBeforeNow(r.Stop)
	if err != nil {
		return fmt.Errorf("range (stop) parse failed: %s", err.Error())
	}

	if !start.Before(stop) {
		return fmt.Errorf("range parse failed: start should before stop")
	}
	return nil
}

// GeneralQuery : the general query option. User can also design their own query option
type GeneralQuery struct {
	Bucket      string    `json:"bucket"`
	Measurement string    `json:"measurement"`
	Fields      []string  `json:"fields"`
	Filters     []kv.KV   `json:"filters"`
	Aggregate   Aggregate `json:"aggregate"`
	Range       Range     `json:"range"`
}

func (g GeneralQuery) Validate() error {
	if g.Bucket == "" {
		return fmt.Errorf("bucket could not be empty")
	}
	if g.Measurement == "" {
		return fmt.Errorf("measurement could not be empty")
	}
	if len(g.Fields) == 0 {
		return fmt.Errorf("must provide at least one field")
	}
	for _, f := range g.Filters {
		if f.Key == "" || f.Value == "" {
			return fmt.Errorf("filter key(value) could not be empty")
		}
	}
	if err := g.Aggregate.Validate(); err != nil {
		return err
	}
	if err := g.Range.Validate(); err != nil {
		return err
	}
	return nil
}

const (
	BucketSnippet      string = "from(bucket: \"%s\")"
	TimeRangeSnippet   string = " |> range(start: %s, stop: %s)"
	FilterSnippet      string = " |> filter(fn: (r) => %s)"
	AggregateSnippet   string = " |> aggregateWindow(every: %s, fn: %s, createEmpty: %v)"
	YieldSnippet       string = " |> yield(name: \"%s\")"
	JoinSnippet        string = "join(tables: {%s}, on: [\"_time\"])"
	MeasurementSnippet string = "r._measurement == \"%s\""
	FieldSnippet       string = "r._field == \"%s\""
	TagSnippet         string = "r.%s == \"%v\""
	GroupSnippet       string = "|> group(columns: [\"sensor_mac\", \"sensor_type\", \"receive_no\"])"
	DropSnippet        string = "|> drop(fn: (column) => column != \"_value\" and column != \"_time\" and column != \"sensor_type\" and column != \"sensor_mac\")"
)

func (g GeneralQuery) TransToFlux() string {
	var scripts []string
	// bucket
	scripts = append(scripts, fmt.Sprintf(BucketSnippet, g.Bucket))

	// range
	scripts = append(scripts, fmt.Sprintf(TimeRangeSnippet, g.Range.Start, g.Range.Stop))

	// measurement
	measurement := fmt.Sprintf("r._measurement == \"%s\"", g.Measurement)
	scripts = append(scripts, fmt.Sprintf(FilterSnippet, measurement))

	// fields
	var fields []string
	for _, f := range g.Fields {
		fields = append(fields, fmt.Sprintf("r._field == \"%s\"", f))
	}
	scripts = append(scripts, fmt.Sprintf(FilterSnippet, strings.Join(fields, " or ")))

	// filters
	if len(g.Filters) > 0 {
		var filters []string
		for _, f := range g.Filters {
			filters = append(filters, fmt.Sprintf("r.%s == \"%v\"", f.Key, f.Value))
		}
		scripts = append(scripts, fmt.Sprintf(FilterSnippet, strings.Join(filters, " and ")))
	}

	// aggregate
	if g.Aggregate.Enable {
		scripts = append(scripts, fmt.Sprintf(AggregateSnippet, g.Aggregate.Every, g.Aggregate.Fn, g.Aggregate.CreateEmpty))
	}

	scripts = append(scripts, GroupSnippet)

	scripts = append(scripts, DropSnippet)

	return strings.Join(scripts, "\n")
}

// CheckTimeBeforeNow : check if the given string is a valid duration string or utc datetime string,
// and check if the time before now.
func CheckTimeBeforeNow(str string) (time.Time, error) {
	now := time.Now()

	// str equals now()
	if str == "now()" {
		return now, nil
	}
	// str is relative time duration, like -20h5m2s
	if d, err := time.ParseDuration(str); err == nil {
		t := now.Add(d)
		if t.Before(now) {
			return t, nil
		}
		return time.Time{}, errors.New("time should before now")
	}
	// str is absolute time format with utc datetime string, like 2021-10-02T15:04:05Z
	if t, err := time.Parse("2006-01-02T15:04:05Z", str); err == nil {
		if t.Before(now) {
			return t, nil
		}
		return time.Time{}, errors.New("time should before now")
	}
	return time.Time{}, errors.New("invalid string")
}

// CheckDurationPositive : check if the given string is a valid duration string,
// and check if the duration is positive. A positive duration like 5h30m, while the negative like -5h30m
func CheckDurationPositive(str string) (time.Duration, error) {
	if d, err := time.ParseDuration(str); err != nil {
		return 0, fmt.Errorf("invaild duration string")
	} else {
		now := time.Now()
		if now.Add(d).Before(now) {
			return 0, fmt.Errorf("negative duration string")
		}
		return d, nil
	}
}
