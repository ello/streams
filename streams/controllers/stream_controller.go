package controllers

import (
	"errors"
	"fmt"
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
	log.Debug("About to register")

	router.POST("/streams", c.handle(c.coalesceStreams))
	router.GET("/stream/:id", c.handle(c.getStream))
	router.POST("/stream/:id", c.handle(c.addToStream))

	log.Debug("Routes Registered %v", router)

}

func (c *streamController) coalesceStreams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	fmt.Println("coalesce")
	return nil
}

func (c *streamController) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	log.WithFields(log.Fields{
		"request": r,
		"params":  ps,
		"id":      ps.ByName("id"),
	}).Debug("/getStream")

	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	log.WithFields(log.Fields{
		"request": r,
		"params":  ps,
		"id":      ps.ByName("id"),
	}).Debug("/addToStream")

	c.JSON(w, http.StatusCreated, "")
	return errors.New("FAIL IT ALREADY")
}
