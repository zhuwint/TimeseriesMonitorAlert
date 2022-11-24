package server

import (
	"sync"

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
	// api := router.Group("/api")
	
}
