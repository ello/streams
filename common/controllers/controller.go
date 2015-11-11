package controllers

import (
	"net/http"

	"github.com/ello/ello-api/controllers"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

//Convienance for the handle function
type action func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) error

//Controller is the interface all of our controllers must implement.
type Controller interface {
	//Register takes a router object and allows the controller to add its routes
	Register(router *httprouter.Router)
}

//This is simply a base struct for purposes of any global storage, to define
//handle off of, and for other controllers to inherit from.  It is not exported.
type baseController struct {
	render.Render
}

//This is a helper function for providing generic error handling for any controllers
//that choose to wrap their actions with it.  It's possible we could redefine this as ServeHTTP
//to avoid the need to call it, but we'll see.
func (c *baseController) handle(a action) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := a(w, r, ps)
		if err != nil {
			switch e := err.(type) {
			// This refers to controllers.Error
			case controllers.Error:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				// log.Printf("HTTP %d - %s", e.Status(), e)
				http.Error(w, e.Error(), e.Status())
			default:
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	})
}
