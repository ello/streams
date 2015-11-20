package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	common "github.com/ello/ello-go/common/http"
	"github.com/ello/ello-go/streams/model"
	"github.com/ello/ello-go/streams/service"
	"github.com/julienschmidt/httprouter"
	"github.com/m4rw3r/uuid"
	"github.com/rcrowley/go-metrics"
)

var addToStreamTimer metrics.Timer
var coalesceTimer metrics.Timer
var getStreamTimer metrics.Timer

type streamController struct {
	baseController
	streamService service.StreamService
}

//NewStreamController is the exported constructor for a streams controller
func NewStreamController(service service.StreamService) Controller {
	addToStreamTimer = metrics.NewTimer()
	coalesceTimer = metrics.NewTimer()
	getStreamTimer = metrics.NewTimer()
	metrics.Register("Streams/AddToStream", addToStreamTimer)
	metrics.Register("Streams/Coalesce", coalesceTimer)
	metrics.Register("Streams/GetStream", getStreamTimer)

	return &streamController{streamService: service}
}

func (c *streamController) Register(router *httprouter.Router) {
	router.PUT("/streams", c.handle(c.addToStream))
	router.POST("/streams/coalesce", c.handle(c.coalesceStreams))
	router.GET("/stream/:id", c.handle(c.getStream))

	log.Debug("Routes Registered")
}

func (c *streamController) coalesceStreams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	startTime := time.Now()
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/coalesce")

	var query model.StreamQuery
	err = json.Unmarshal(body, &query)

	if err != nil {
		return common.StatusError{Code: 422, Err: err}
	}
	items, err := c.streamService.Load(query)
	if err != nil {
		return common.StatusError{Code: 400, Err: errors.New("An error occurred loading streams")}
	}

	c.JSON(w, http.StatusOK, items)
	coalesceTimer.UpdateSince(startTime)
	return nil
}

func (c *streamController) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	startTime := time.Now()
	log.WithFields(fieldsFor(r, nil, nil)).Debug("/getStream")

	//get ID and validate that it is a uuid.
	streamID, err := uuid.FromString(ps.ByName("id"))
	if err != nil && !streamID.IsZero() {
		return common.StatusError{Code: 422, Err: errors.New("id must be a valid UUID")}
	}
	items, err := c.streamService.Load(model.StreamQuery{Streams: []uuid.UUID{streamID}})
	if err != nil {
		return common.StatusError{Code: 400, Err: errors.New("An error occurred loading streams")}
	}

	c.JSON(w, http.StatusOK, items)
	getStreamTimer.UpdateSince(startTime)
	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	startTime := time.Now()
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/addToStream")

	var items []model.StreamItem
	err = json.Unmarshal(body, &items)

	log.WithFields(log.Fields{
		"items": items,
		"err":   err,
	}).Debug("Unmarshaled items")

	if err != nil {
		return common.StatusError{Code: 422, Err: errors.New("body must be an array of StreamItems")}
	}

	err = c.streamService.Add(items)

	if err != nil {
		return common.StatusError{Code: 400, Err: errors.New("An error occurred adding to the stream(s)")}
	}

	c.JSON(w, http.StatusCreated, nil)
	addToStreamTimer.UpdateSince(startTime)
	return nil
}
