package controllers

import (
	"fmt"
	"net/http"

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
	router.POST("/streams", c.Handle(c.coalesceStreams))
	router.GET("/stream/:id", c.Handle(c.getStream))
	router.POST("/stream/:id", c.Handle(c.addToStream))
}

func (c *streamController) coalesceStreams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	fmt.Println("coalesce")
	return nil
}

func (c *streamController) getStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	fmt.Println("getStream")
	return nil
}

func (c *streamController) addToStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	fmt.Println("addToStream")
	return nil
}
