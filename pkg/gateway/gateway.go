package gateway

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"timeseries/pkg/api"
	"timeseries/pkg/models"

	"github.com/gin-gonic/gin"
)

type server struct {
	httpMux *gin.Engine
	port    int

	// point process queue
	queue chan models.Point
	// stop signal
	stopH chan struct{}
}

func NewServer(port, buffSize int) *server {
	return &server{
		httpMux: gin.New(),
		port:    port,
		queue:   make(chan models.Point, buffSize),
		stopH:   make(chan struct{}),
	}
}

func (s *server) Start() error {
	s.registerRoutes()

	go s.process()

	if err := s.httpMux.Run(fmt.Sprintf(":%d", s.port)); err != nil {
		logrus.Errorf("http server exist with error: %s", err.Error())
	}
	return nil
}

func (s *server) Stop() error {
	// close stop channel for send signal to all goroutines
	close(s.stopH)
	return nil
}

func (s *server) registerRoutes() {
	s.httpMux.Handle(http.MethodGet, "/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
	})

	apiRouteV2 := s.httpMux.Group("/api/v2")
	{
		apiRouteV2.Handle(http.MethodPost, "/write", s.pointReceiver)
	}
}

func (s *server) pointReceiver(ctx *gin.Context) {
	precision := ctx.Query("precision")
	if precision == "" {
		precision = "n"
	}
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	points, err := models.ParsePointsWithPrecision(body, time.Now().UTC(), precision)
	if err != nil {
		if err == io.EOF {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusBadRequest, api.ReplyError{Code: api.RequestBodyError, Msg: err.Error()})
		return
	}
	// response with async
	go s.write(points)
	ctx.JSON(http.StatusNoContent, nil)
}

func (s *server) write(points []models.Point) {
	for i := range points {
		s.queue <- points[i]
	}
}

func (s *server) process() {
	for {
		select {
		case p := <-s.queue:
			s.storeToInfluxdb(p)
			s.publish(p)
		case <-s.stopH:
			logrus.Infof("stop point process")
			break
		}
	}
}

// write point to influxdb
func (s *server) storeToInfluxdb(p models.Point) {

}

// publish point to dapr
func (s *server) publish(p models.Point) {

}
