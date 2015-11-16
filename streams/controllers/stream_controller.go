package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	common "github.com/ello/ello-go/common/http"
	"github.com/julienschmidt/httprouter"
	"github.com/m4rw3r/uuid"
)

type streamController struct {
	baseController
}

// TODO Move all of the stream item stuff to its own file, probably in a data pkg
type streamItemType int

const (
	//POST is a type of stream item which is a direct post
	POST streamItemType = iota
	//REPOST is a type of stream item which represents a repost
	REPOST
)

//StreamItem represents a single item on a stream
type StreamItem struct {
	ID        uuid.UUID      `json:"id"`
	Timestamp time.Time      `json:"ts"`
	Type      streamItemType `json:"type"`
	StreamID  uuid.UUID      `json:"stream_id"`
}

//NewStreamController is the exported constructor for a streams controller
func NewStreamController() Controller {
	return &streamController{}
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

	streamID, _ := uuid.V4()
	items := generateFakeResponse(streamID)
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
	items := generateFakeResponse(streamID)
	c.JSON(w, http.StatusOK, items)

	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/addToStream")

	// do stuff
	var items []StreamItem
	err = json.Unmarshal(body, &items)

	log.WithFields(log.Fields{
		"items": items,
		"err":   err,
	}).Debug("Unmarshaled items")

	if err != nil {
		return common.StatusError{Code: 422, Err: errors.New("body must be an array of StreamItems")}
	}

	c.JSON(w, http.StatusCreated, nil)
	return nil
}

func generateFakeResponse(streamID uuid.UUID) []StreamItem {
	//fake data
	uuid1, _ := uuid.V4()
	uuid2, _ := uuid.V4()

	return []StreamItem{
		{
			ID:        uuid1,
			Timestamp: time.Now(),
			Type:      POST,
			StreamID:  streamID,
		}, {
			ID:        uuid2,
			Timestamp: time.Now(),
			Type:      REPOST,
			StreamID:  streamID,
		},
	}
}
