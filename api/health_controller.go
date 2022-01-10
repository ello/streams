package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

type healthController struct {
	baseController
	startTime time.Time
	commit    string
	roshi     string
}

type heartbeatResponse struct {
	Commit string `json:"commit"`
	Uptime string `json:"uptime"`
}

//NewHealthController returns a controller that will display metrics to /metrics
func NewHealthController(startTime time.Time, commit string, roshiURI string) Controller {
	return &healthController{
		startTime: startTime,
		commit:    commit,
		roshi:     roshiURI,
	}
}

func (c *healthController) Register(router *httprouter.Router) {
	router.GET("/health/metrics", c.handle(c.printMetrics))
	router.GET("/health/check", c.handle(c.healthCheck))
	router.GET("/health/heartbeat", c.handle(c.heartbeat))

	log.Debug("Health Routes Registered")
}

// printMetrics will print all metrics from the default registry in the response
func (c *healthController) printMetrics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	buffer := new(bytes.Buffer)
	metrics.WriteOnce(metrics.DefaultRegistry, buffer)
	c.Text(w, http.StatusOK, buffer.String())
	return nil
}

// healthCheck will verify it can communicate to the configured roshi instance and return ERR/OK appropriately
func (c *healthController) healthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(c.roshi + "/metrics")
	if err != nil || resp.StatusCode != 200 {
		c.Text(w, http.StatusInternalServerError, "ERR")
		return nil
	}
	c.Text(w, http.StatusOK, "OK")

	return nil
}

// heartbeat will return a response with the uptime and commit the binary was built with (if available)
func (c *healthController) heartbeat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	heartbeat := heartbeatResponse{
		Commit: c.commit,
		Uptime: time.Since(c.startTime).String(),
	}

	c.JSON(w, http.StatusOK, heartbeat)
	return nil
}
