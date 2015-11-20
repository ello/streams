package api

import (
	"bytes"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/rcrowley/go-metrics"
)

type metricsController struct {
	baseController
}

//NewMetricsController returns a controller that will display metrics to /metrics
func NewMetricsController() Controller {
	return &metricsController{}
}

func (c *metricsController) Register(router *httprouter.Router) {
	router.GET("/metrics", c.handle(c.printMetrics))

	log.Debug("Metrics Routes Registered")
}

func (c *metricsController) printMetrics(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	buffer := new(bytes.Buffer)
	metrics.WriteOnce(metrics.DefaultRegistry, buffer)
	c.Text(w, http.StatusOK, buffer.String())
	return nil
}
