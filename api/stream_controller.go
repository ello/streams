package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/streams/model"
	"github.com/ello/streams/service"
	"github.com/ello/streams/util"
	"github.com/julienschmidt/httprouter"
	"github.com/rcrowley/go-metrics"
)

var addToStreamTimer metrics.Timer
var removeFromStreamTimer metrics.Timer
var coalesceTimer metrics.Timer
var getStreamTimer metrics.Timer

type streamController struct {
	baseController
	streamService service.StreamService
	authConfig    AuthConfig
}

//NewStreamController is the exported constructor for a streams controller
func NewStreamController(service service.StreamService, authConfig AuthConfig) Controller {
	addToStreamTimer = metrics.NewTimer()
	removeFromStreamTimer = metrics.NewTimer()
	coalesceTimer = metrics.NewTimer()
	getStreamTimer = metrics.NewTimer()
	metrics.Register("Streams.AddToStream", addToStreamTimer)
	metrics.Register("Streams.RemoveFromStream", removeFromStreamTimer)
	metrics.Register("Streams.Coalesce", coalesceTimer)
	metrics.Register("Streams.GetStream", getStreamTimer)

	return &streamController{streamService: service, authConfig: authConfig}
}

func (c *streamController) Register(router *httprouter.Router) {
	router.PUT("/streams", basicAuth(timeRequest(c.handle(c.addToStream), addToStreamTimer), c.authConfig))
	router.DELETE("/streams", basicAuth(timeRequest(c.handle(c.removeFromStream), removeFromStreamTimer), c.authConfig))
	router.POST("/streams/coalesce", basicAuth(timeRequest(c.handle(c.coalesceStreams), coalesceTimer), c.authConfig))
	router.GET("/stream/:id", basicAuth(timeRequest(c.handle(c.getStream), getStreamTimer), c.authConfig))

	log.Debug("Routes Registered")
}

func (c *streamController) coalesceStreams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/coalesce")

	queryParams := r.URL.Query()
	limit, err := util.ValidateInt(queryParams.Get("limit"), 10)
	if err != nil {
		return StatusError{Code: 422, Err: errors.New("Limit should be a number")}
	}

	fromSlug := queryParams.Get("from")

	var query model.StreamQuery
	err = json.Unmarshal(body, &query)

	if err != nil {
		return StatusError{Code: 422, Err: err}
	}

	response, err := c.streamService.Load(query, limit, fromSlug)
	if err != nil {
		return StatusError{Code: 400, Err: errors.New("An error occurred loading streams")}
	}

	addLink(w, nextPage(r, response, limit))

	c.JSON(w, http.StatusOK, response.Items)
	return nil
}

func (c *streamController) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	log.WithFields(fieldsFor(r, nil, nil)).Debug("/getStream")

	//get ID and validate that it is a uuid.
	streamID := ps.ByName("id")

	queryParams := r.URL.Query()
	limit, err := util.ValidateInt(queryParams.Get("limit"), 10)
	if err != nil {
		return StatusError{Code: 422, Err: errors.New("Limit should be a number")}
	}
	fromSlug := queryParams.Get("from")

	response, err := c.streamService.Load(model.StreamQuery{Streams: []string{streamID}}, limit, fromSlug)
	if err != nil {
		return StatusError{Code: 400, Err: errors.New("An error occurred loading streams")}
	}

	addLink(w, nextPage(r, response, limit))

	c.JSON(w, http.StatusOK, response.Items)
	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	return c.performStreamAction(w, r, c.streamService.Add, http.StatusCreated)
}

func (c *streamController) removeFromStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	return c.performStreamAction(w, r, c.streamService.Remove, http.StatusOK)
}

func (c *streamController) performStreamAction(w http.ResponseWriter, r *http.Request, action func([]model.StreamItem) error, status int) error {
	items, err := getItemsFromBody(r, "/updateStream")

	if err != nil {
		return err
	}

	err = action(items)

	if err != nil {
		return StatusError{Code: 400, Err: errors.New("An error occurred removing from the stream(s)")}
	}

	c.JSON(w, status, nil)
	return nil
}

func addLink(w http.ResponseWriter, nextPageLink string) {
	w.Header().Set("Link", fmt.Sprintf("<%v>; rel=\"next\"", nextPageLink))
}

func nextPage(r *http.Request, items *model.StreamQueryResponse, limit int) string {
	uri := ""
	if r.TLS != nil {
		uri = "https://"
	} else {
		uri = "http://"
	}

	return fmt.Sprintf("%v%v%v?limit=%d&from=%s", uri, r.Host, r.URL.Path, limit, items.Cursor)
}

func getItemsFromBody(r *http.Request, debugName string) ([]model.StreamItem, error) {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug(debugName)

	var items []model.StreamItem
	err = json.Unmarshal(body, &items)

	log.WithFields(log.Fields{
		"items": items,
		"err":   err,
	}).Debug("Unmarshaled items")

	if err != nil {
		return items, StatusError{Code: 422, Err: errors.New("body must be an array of StreamItems")}
	}

	return items, err
}

func timeRequest(action httprouter.Handle, timer metrics.Timer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		startTime := time.Now()
		action(w, r, ps)
		timer.UpdateSince(startTime)
		return
	}
}
