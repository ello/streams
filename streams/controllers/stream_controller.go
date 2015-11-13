package controllers

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type streamController struct {
	baseController
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

	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	log.WithFields(fieldsFor(r, body, err)).Debug("/addToStream")

	c.JSON(w, http.StatusCreated, "")
	return nil
}

func fieldsFor(r *http.Request, body []byte, err error) log.Fields {
	return log.Fields{
		"url":     r.URL,
		"method":  r.Method,
		"headers": r.Header,
		"body":    string(body[:]),
		"err":     err,
	}
}
