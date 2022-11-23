package server

import (
	"sync"

	"timeseries/cmd/data-manager/server/sensor"
	"timeseries/cmd/data-manager/server/timeseries"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	router *gin.Engine
	once   = &sync.Once{}
)

func GetRouter() *gin.Engine {
	once.Do(func() {
		router = gin.New()
		router.Use(gin.Recovery())
		initRoute()
	})

	defer func() {
		if err := recover(); err != nil {
			logrus.Error("server runtime failed")
		}
	}()

	return router
}

func initRoute() {
	api := router.Group("/api")
	{
		api.GET("/project", sensor.GetProjects)
		api.GET("/sensor", sensor.GetSensors)
		api.GET("/location", sensor.GetLocations)
		api.GET("/measurement", sensor.GetMeasurements)
	}
	{
		api.POST("/data", timeseries.QueryTimeseries)
	}
}
