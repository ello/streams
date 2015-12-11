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
	librato "github.com/mihasya/go-metrics-librato"
	metrics "github.com/rcrowley/go-metrics"
)

var commit string
var startTime = time.Now()
var help bool
var helpMessage = `ELLO STREAM API
--------------------------
Set ENV Variables to configure:

ELLO_API_PORT for the port to run this service on.  Default is 8080
ELLO_ROSHI_HOST for the location of the roshi instance.  Default is http://localhost:6302
ELLO_ROSHI_TIMEOUT for the timeout (in Seconds) for roshi connections.  Default is 5s.
ELLO_AUTH_ENABLED any value will enable basic auth.  Default is disabled.
ELLO_AUTH_USERNAME for the auth username.  Default is 'ello'.
ELLO_AUTH_PASSWORD for the auth password.  Default is 'password'.
ELLO_LOG_LEVEL for the log level.  Valid levels are "debug", "info", "warn", "error".  Default is warn. 
ELLO_LIBRATO_EMAIL librato config
ELLO_LIBRATO_TOKEN librato config
ELLO_LIBRATO_HOSTNAME librato config
`

func main() {

	flag.BoolVar(&help, "h", false, "help?")
	flag.Parse()

	if help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	var logLevel log.Level

	switch util.GetEnvWithDefault("ELLO_LOG_LEVEL", "warn") {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "error":
		logLevel = log.ErrorLevel
	default:
		logLevel = log.WarnLevel
	}

	log.SetLevel(logLevel)
	fmt.Printf("Using log level [%v]\n", logLevel)

	roshi := util.GetEnvWithDefault("ELLO_ROSHI_HOST", "http://localhost:6302")
	streamsService, err := service.NewRoshiStreamService(roshi, time.Duration(util.GetEnvIntWithDefault("ELLO_ROSHI_TIMEOUT", 5))*time.Second)
	if err != nil {
		log.Panic(err)
	}

	authConfig := api.AuthConfig{
		Username: []byte(util.GetEnvWithDefault("ELLO_AUTH_USERNAME", "ello")),
		Password: []byte(util.GetEnvWithDefault("ELLO_AUTH_PASSWORD", "password")),
		Enabled:  util.IsEnvPresent("ELLO_AUTH_ENABLED"),
	}
	log.Infof(authConfig.String())

	if util.IsEnvPresent("ELLO_LIBRATO_TOKEN") {
		go librato.Librato(metrics.DefaultRegistry,
			10e9, // interval
			os.Getenv("ELLO_LIBRATO_EMAIL"),    // account owner email address
			os.Getenv("ELLO_LIBRATO_TOKEN"),    // Librato API token
			os.Getenv("ELLO_LIBRATO_HOSTNAME"), // source
			[]float64{0.95},                    // percentiles to send
			time.Millisecond,                   // time unit
		)
	}

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
