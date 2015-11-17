package main

import (
	"net/http"

	"github.com/ello/ello-go/common/util"
	"github.com/ello/ello-go/streams/api"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	streamsController := api.NewStreamController()

	// controllers register their routes with the router
	streamsController.Register(router)

	port := util.GetEnvWithDefault("ELLO_API_PORT", "8080")
	http.ListenAndServe(":"+port, router)
}
