package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	common "github.com/ello/ello-go/common/http"
	"github.com/ello/ello-go/streams/model"
	"github.com/ello/ello-go/streams/service"
	"github.com/julienschmidt/httprouter"
	"github.com/m4rw3r/uuid"
)

type streamController struct {
	baseController
	service service.StreamService
}

//NewStreamController is the exported constructor for a streams controller
func NewStreamController(service service.StreamService) Controller {
	return &streamController{service: service}
}

func (c *streamController) Register(router *httprouter.Router) {
	router.PUT("/streams", c.handle(c.addToStream))
	router.POST("/streams/coalesce", c.handle(c.coalesceStreams))
	router.GET("/stream/:id", c.handle(c.getStream))

	log.Debug("Routes Registered")
}

func (c *streamController) coalesceStreams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/coalesce")

	var query model.StreamQuery
	err = json.Unmarshal(body, &query)

	if err != nil {
		return common.StatusError{Code: 422, Err: err}
	}
	items, err := c.service.LoadContent(query)
	if err != nil {
		return common.StatusError{Code: 400, Err: errors.New("An error occurred loading streams")}
	}

	c.JSON(w, http.StatusOK, items)
	return nil
}

func (c *streamController) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	log.WithFields(fieldsFor(r, nil, nil)).Debug("/getStream")

	//get ID and validate that it is a uuid.
	streamID, err := uuid.FromString(ps.ByName("id"))
	if err != nil && !streamID.IsZero() {
		return common.StatusError{Code: 422, Err: errors.New("id must be a valid UUID")}
	}
	items, err := c.service.LoadContent(model.StreamQuery{Streams: []uuid.UUID{streamID}})
	if err != nil {
		return common.StatusError{Code: 400, Err: errors.New("An error occurred loading streams")}
	}

	c.JSON(w, http.StatusOK, items)
	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
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

	err = c.service.AddContent(items)

	if err != nil {
		return common.StatusError{Code: 400, Err: errors.New("An error occurred adding to the stream(s)")}
	}

	c.JSON(w, http.StatusCreated, nil)
	return nil
}
