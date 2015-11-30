package api

import (
	"bytes"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/rcrowley/go-metrics"
)

type healthController struct {
	baseController
	startTime time.Time
	commit    string
}

type heartbeat struct {
	Commit string `json:"commit"`
	Uptime string `json:"uptime"`
}

//NewHealthController returns a controller that will display metrics to /metrics
func NewHealthController(startTime time.Time, commit string) Controller {
	return &healthController{
		startTime: startTime,
		commit:    commit,
	}
}

func (c *healthController) Register(router *httprouter.Router) {
	router.GET("/health/metrics", c.handle(c.printMetrics))
	router.GET("/health/check", c.handle(c.healthCheck))
	router.GET("/health/heartbeat", c.handle(c.heartbeat))

	log.Debug("Health Routes Registered")
}

func (c *healthController) printMetrics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	//TODO see if metrics supports json output
	buffer := new(bytes.Buffer)
	metrics.WriteOnce(metrics.DefaultRegistry, buffer)
	c.Text(w, http.StatusOK, buffer.String())
	return nil
}

func (c *healthController) healthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	//TODO Make this more robust
	c.Text(w, http.StatusOK, "OK")
	return nil
}

func (c *healthController) heartbeat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	heartbeat := heartbeat{
		Commit: c.commit,
		Uptime: time.Now().Sub(c.startTime).String(),
	}

	c.JSON(w, http.StatusOK, heartbeat)
	return nil
}
