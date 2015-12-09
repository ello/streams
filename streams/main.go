package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/ello/ello-go/streams/api"
	"github.com/ello/ello-go/streams/service"
	"github.com/ello/ello-go/streams/util"
	"github.com/julienschmidt/httprouter"
	nlog "github.com/meatballhat/negroni-logrus"
	"github.com/rcrowley/go-metrics"
)

var commit string
var startTime = time.Now()
var verbose bool
var veryVerbose bool
var help bool
var helpMessage = `ELLO STREAM API
--------------------------
-v or -vv for verbose/veryverbose output

Set ENV Variables to configure:
ELLO_API_PORT for the port to run this service on.  Default is 8080
ELLO_ROSHI_HOST for the location of the roshi instance.  Default is http://localhost:6302
ELLO_AUTH_ENABLED any value will enable basic auth.  Default is disabled.
ELLO_AUTH_USERNAME for the auth username.  Default is 'ello'.
ELLO_AUTH_PASSWORD for the auth password.  Default is 'password'.
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
		go metrics.Write(metrics.DefaultRegistry, 1*time.Minute, log.New().Writer())
	}
	log.SetLevel(logLevel)

	roshi := util.GetEnvWithDefault("ELLO_ROSHI_HOST", "http://localhost:6302")
	streamsService, err := service.NewRoshiStreamService(roshi)
	if err != nil {
		log.Panic(err)
	}

	authConfig := api.AuthConfig{
		Username: []byte(util.GetEnvWithDefault("ELLO_AUTH_USERNAME", "ello")),
		Password: []byte(util.GetEnvWithDefault("ELLO_AUTH_PASSWORD", "password")),
		Enabled:  util.IsEnvPresent("ELLO_AUTH_ENABLED"),
	}
	log.Infof(authConfig.String())

	router := httprouter.New()

	streamsController := api.NewStreamController(streamsService, authConfig)
	streamsController.Register(router)

	healthController := api.NewHealthController(startTime, commit, roshi)
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
