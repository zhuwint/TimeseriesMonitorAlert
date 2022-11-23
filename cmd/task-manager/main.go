/*
Copyright © 2022 zhuwentao
*/
package main

import (
	"flag"

	"timeseries/pkg/utils/env"
	"timeseries/pkg/vars"
)

var (
	// mysql config
	MysqlAddress     = flag.String(vars.MYSQL_ADDRESS, "localhost:3306", "mysql server address")
	MysqlUser  = flag.String(vars.MYSQL_USER, "", "mysql user")
	MysqlPassword = flag.String(vars.MYSQL_PASSWORD, "", "mysql password")
	// influxdb config
	InfluxAddress  = flag.String(vars.INFLUX_ADDRESS, "localhost", "influx server address")
	InfluxToken = flag.String(vars.INFLUX_TOKEN, "", "influx token")
	InfluxOrg   = flag.String(vars.INFLUX_ORG, "", "influx org")
	// service config
	ServiceName = flag.String(vars.SERVICE_NAME, "task-manager", "service name")
	ServicePort = flag.Int(vars.SERVICE_PORT, 3000, "service port")
)

func main() {
	flag.Parse()

	*InfluxAddress = env.GetEnvString(vars.INFLUX_ADDRESS, *InfluxAddress)
	*InfluxToken = env.GetEnvString(vars.INFLUX_TOKEN, *InfluxToken)
	*InfluxOrg = env.GetEnvString(vars.INFLUX_ORG, *InfluxOrg)

}
