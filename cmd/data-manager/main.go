/*
Copyright © 2022 zhuwentao
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"timeseries/cmd/data-manager/server"
	influxsvc "timeseries/pkg/service/influxdb"
	mysqlsvc "timeseries/pkg/service/mysql"
	"timeseries/pkg/utils/env"
	"timeseries/pkg/vars"

	"github.com/sirupsen/logrus"
)

var (
	// mysql config
	MysqlAddress  = flag.String(vars.MYSQL_ADDRESS, "localhost:3306", "mysql server address")
	MysqlUser     = flag.String(vars.MYSQL_USER, "", "mysql user")
	MysqlPassword = flag.String(vars.MYSQL_PASSWORD, "", "mysql password")
	MysqlDatabase = flag.String(vars.MYSQL_DATABASE, "", "mysql database")
	// influxdb config
	InfluxAddress = flag.String(vars.INFLUX_ADDRESS, "localhost", "influx server address")
	InfluxToken   = flag.String(vars.INFLUX_TOKEN, "", "influx token")
	InfluxOrg     = flag.String(vars.INFLUX_ORG, "", "influx org")
	InfluxBucket  = flag.String(vars.INFLUX_BUCKET, "", "influx bucket")
	// service config
	ServiceName = flag.String(vars.SERVICE_NAME, "data-manager", "service name")
)

func main() {
	flag.Parse()

	parseEnvs()

	if err := initMysqlService(); err != nil {
		logrus.Error("init mysql service failed:", err.Error())
		return
	}

	if err := initInfluxService(); err != nil {
		logrus.Error("init influxdb service failed:", err.Error())
		return
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 3000),
		Handler: server.GetRouter(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("service start failed: %s", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Error("Server Shutdown:", err)
	}

	// user for control service close
	// TODO: close goroutine and service, save state
	exit, exitCancel := context.WithCancel(context.Background())
	defer exitCancel()

	select {
	case <-ctx.Done():
		logrus.Info("timeout of 3 seconds.")
	case <-exit.Done():
		influxsvc.CloseClient()
		logrus.Info("all service closed")
	}
	logrus.Info("Server exiting")
	os.Exit(0)
}

func parseEnvs() {
	*InfluxAddress = env.GetEnvString(vars.INFLUX_ADDRESS, *InfluxAddress)
	*InfluxToken = env.GetEnvString(vars.INFLUX_TOKEN, *InfluxToken)
	*InfluxOrg = env.GetEnvString(vars.INFLUX_ORG, *InfluxOrg)
	*InfluxBucket = env.GetEnvString(vars.INFLUX_BUCKET, *InfluxBucket)

	*MysqlAddress = env.GetEnvString(vars.MYSQL_ADDRESS, *MysqlAddress)
	*MysqlUser = env.GetEnvString(vars.MYSQL_USER, *MysqlUser)
	*MysqlPassword = env.GetEnvString(vars.MYSQL_PASSWORD, *MysqlPassword)
	*MysqlDatabase = env.GetEnvString(vars.MYSQL_DATABASE, *MysqlDatabase)
}

func initMysqlService() error {
	mysqlAccount := mysqlsvc.Account{
		Address:  *MysqlAddress,
		Database: *MysqlDatabase,
		Username: *MysqlUser,
		Password: *MysqlPassword,
	}

	if err := mysqlAccount.Validate(); err != nil {
		return err
	}

	mysqlsvc.InitMysqlClient(mysqlAccount)
	logrus.Infof("init mysql success using config database:%s username:%s url:%s", *MysqlDatabase, *MysqlUser, *MysqlAddress)
	return nil
}

func initInfluxService() error {
	influxAccount := influxsvc.Account{
		Address: *InfluxAddress,
		Org:     *InfluxOrg,
		Token:   *InfluxToken,
		Bucket:  *InfluxBucket,
	}

	if err := influxAccount.Validate(); err != nil {
		return err
	}

	influxsvc.InitInfluxClient(influxAccount, *InfluxBucket)
	logrus.Infof("init influx sucess using config bucket:%s org:%s url:%s", *InfluxBucket, *InfluxOrg, *InfluxAddress)
	return nil
}
