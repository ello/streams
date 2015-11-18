package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/ello-go/common/util"
	"github.com/ello/ello-go/streams/api"
	"github.com/ello/ello-go/streams/service"
	"github.com/julienschmidt/httprouter"
)

func main() {
	log.SetLevel(log.DebugLevel)

	streamsService, err := service.NewRoshiStreamService("http://localhost:6302")
	if err != nil {
		log.Panic(err)
	}

	router := httprouter.New()

	router.GET("/test", httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Debug("test")
		w.Write([]byte("test ack!"))
	}))

	streamsController := api.NewStreamController(streamsService)

	// controllers register their routes with the router
	streamsController.Register(router)

	port := util.GetEnvWithDefault("ELLO_API_PORT", "8080")
	serverAt := "localhost:" + port
	log.Debugf("Listening at: %v", serverAt)
	http.ListenAndServe(serverAt, router)
}
