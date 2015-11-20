package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/ello/ello-go/common/util"
	"github.com/ello/ello-go/streams/api"
	"github.com/ello/ello-go/streams/service"
	"github.com/julienschmidt/httprouter"
	nlog "github.com/meatballhat/negroni-logrus"
	"github.com/rcrowley/go-metrics"
)

var verbose bool
var veryVerbose bool
var help bool
var helpMessage = `ELLO STREAM API
--------------------------
-v or -vv for verbose/veryverbose output

Set ENV Variables to configure:
ELLO_API_PORT for the port to run this service on.  Default is 8080
ELLO_ROSHI_HOST for the location of the roshi instance.  Default is http://localhost:6302
`

func main() {

	flag.BoolVar(&verbose, "v", false, "if set, show messages from the logger")
	flag.BoolVar(&veryVerbose, "vv", false, "if set, show ALL messages from the logger")
	flag.BoolVar(&help, "h", false, "help?")
	flag.Parse()

	if help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	logLevel := log.WarnLevel
	if verbose {
		logLevel = log.InfoLevel
	}
	if veryVerbose {
		logLevel = log.DebugLevel
		go metrics.Write(metrics.DefaultRegistry, 1*time.Second, log.New().Writer())
	}
	log.SetLevel(logLevel)

	streamsService, err := service.NewRoshiStreamService(util.GetEnvWithDefault("ELLO_ROSHI_HOST", "http://localhost:6302"))
	if err != nil {
		log.Panic(err)
	}

	router := httprouter.New()

	streamsController := api.NewStreamController(streamsService)
	// controllers register their routes with the router
	streamsController.Register(router)

	healthController := api.NewHealthController()
	healthController.Register(router)

	n := negroni.New(
		negroni.NewRecovery(),
		nlog.NewCustomMiddleware(logLevel, &log.TextFormatter{}, "web"),
	)
	n.UseHandler(router)

	port := util.GetEnvWithDefault("ELLO_API_PORT", "8080")
	serverAt := ":" + port
	n.Run(serverAt)
}
