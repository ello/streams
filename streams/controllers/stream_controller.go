package controllers

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/m4rw3r/uuid"
)

type streamController struct {
	baseController
}

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
}

//NewStreamController is the exported constructor for a streams controller
func NewStreamController() Controller {
	return &streamController{}
}

func (c *streamController) Register(router *httprouter.Router) {
	router.POST("/streams", c.handle(c.coalesceStreams))
	router.GET("/stream/:id", c.handle(c.getStream))
	router.POST("/stream/:id", c.handle(c.addToStream))

	log.Debug("Routes Registered")
}

func (c *streamController) coalesceStreams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/coalesce")

	return nil
}

func (c *streamController) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	log.WithFields(fieldsFor(r, nil, nil)).Debug("/getStream")

	//fake data
	uuid1, _ := uuid.V4()
	uuid2, _ := uuid.V4()

	items := []StreamItem{
		{
			ID:        uuid1,
			Timestamp: time.Now(),
			Type:      POST,
		}, {
			ID:        uuid2,
			Timestamp: time.Now(),
			Type:      REPOST,
		},
	}

	c.JSON(w, http.StatusOK, items)

	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/addToStream")

	// do stuff

	c.JSON(w, http.StatusCreated, nil)
	return nil
}
