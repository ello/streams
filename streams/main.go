package main

import (
	"flag"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/ello/ello-go/common/util"
	"github.com/ello/ello-go/streams/api"
	"github.com/ello/ello-go/streams/service"
	"github.com/julienschmidt/httprouter"
	nlog "github.com/meatballhat/negroni-logrus"
)

var verbose bool
var veryVerbose bool

func main() {

	flag.BoolVar(&verbose, "v", false, "if set, show messages from the logger")
	flag.BoolVar(&veryVerbose, "vv", false, "if set, show ALL messages from the logger")

	flag.Parse()

	logLevel := log.WarnLevel
	if verbose {
		logLevel = log.InfoLevel
	}
	if veryVerbose {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

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

	n := negroni.New(
		negroni.NewRecovery(),
		nlog.NewCustomMiddleware(logLevel, &log.TextFormatter{}, "web"),
	)
	n.UseHandler(router)

	port := util.GetEnvWithDefault("ELLO_API_PORT", "8080")
	serverAt := ":" + port
	n.Run(serverAt)
}
