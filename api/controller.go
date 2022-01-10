package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

//Action is a convienance for the handle function
type action func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) error

//Controller is the interface all of our controllers must implement.
type Controller interface {
	//Register takes a router object and allows the controller to add its routes
	Register(router *httprouter.Router)
}

// AuthConfig contains all the necessary configuration for basic auth setup.
type AuthConfig struct {
	Enabled  bool
	Username []byte
	Password []byte
}

//StatusString status string for auth config
func (a AuthConfig) String() string {
	if a.Enabled {
		return fmt.Sprintf("Authentication is Enabled with Username %v and Password %v", string(a.Username), string(a.Password))
	}
	return "Authentication is Disabled"
}

// BaseController is simply a base struct for purposes of any global storage, to define
//handle off of, and for other controllers to inherit from.  It is not exported.
type baseController struct {
	render.Render
}

//Handle is a helper function for providing generic error handling for any controllers
//that choose to wrap their actions with it.
func (c *baseController) handle(a action) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := a(w, r, ps)
		if err != nil {
			switch e := err.(type) {
			// This refers to controllers.Error
			case Error:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				log.Debugf("HTTP %d - %s", e.Status(), e)
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

//basicAuth is a wrapper for a handler to configure basic auth, per the passed AuthConfig
func basicAuth(h httprouter.Handle, authConfig AuthConfig) httprouter.Handle {
	//is auth is disabled, skip wrapping it
	if !authConfig.Enabled {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "

		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 &&
					bytes.Equal(pair[0], authConfig.Username) &&
					bytes.Equal(pair[1], authConfig.Password) {

					// Delegate request to the given handle
					h(w, r, ps)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

func fieldsFor(r *http.Request, body []byte, err error) logrus.Fields {
	return logrus.Fields{
		"url":     r.URL,
		"method":  r.Method,
		"headers": r.Header,
		"body":    string(body[:]),
		"err":     err,
	}
}
