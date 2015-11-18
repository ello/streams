package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/ello-go/common/util"
	"github.com/ello/ello-go/streams/api"
	"github.com/ello/ello-go/streams/service"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	streamsService, err := service.NewRoshiStreamService("http://localhost:6302")
	if err != nil {
		log.Panic(err)
		panic(1)
	}
	streamsController := api.NewStreamController(streamsService)

	// controllers register their routes with the router
	streamsController.Register(router)

	port := util.GetEnvWithDefault("ELLO_API_PORT", "8080")
	fmt.Println("Listening on port: " + port)
	http.ListenAndServe(":"+port, router)
}
