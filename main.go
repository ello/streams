package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/ello/streams/api"
	"github.com/ello/streams/service"
	"github.com/ello/streams/util"
	"github.com/julienschmidt/httprouter"
	nlog "github.com/meatballhat/negroni-logrus"
	librato "github.com/mihasya/go-metrics-librato"
	metrics "github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
)

var commit string
var startTime = time.Now()
var help bool
var helpMessage = `ELLO STREAM API
--------------------------
Set ENV Variables to configure:

PORT for the port to run this service on.  Default is 8080
ROSHI_URL for the location of the roshi instance.  Default is http://localhost:6302
ROSHI_TIMEOUT for the timeout (in Seconds) for roshi connections.  Default is 5s.
AUTH_ENABLED any value will enable basic auth.  Default is disabled.
AUTH_USERNAME for the auth username.  Default is 'ello'.
AUTH_PASSWORD for the auth password.  Default is 'password'.
LOG_LEVEL for the log level.  Valid levels are "debug", "info", "warn", "error".  Default is warn.
LIBRATO_EMAIL librato config
LIBRATO_TOKEN librato config
LIBRATO_HOSTNAME librato config
`

func main() {

	flag.BoolVar(&help, "h", false, "help?")
	flag.Parse()

	if help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	var logLevel log.Level

	switch util.GetEnvWithDefault("LOG_LEVEL", "warn") {
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

	roshi := util.GetEnvWithDefault("ROSHI_URL", "http://localhost:6302")
	streamsService, err := service.NewRoshiStreamService(roshi, time.Duration(util.GetEnvIntWithDefault("ROSHI_TIMEOUT", 5))*time.Second)
	if err != nil {
		log.Panic(err)
	}

	authConfig := api.AuthConfig{
		Username: []byte(util.GetEnvWithDefault("AUTH_USERNAME", "ello")),
		Password: []byte(util.GetEnvWithDefault("AUTH_PASSWORD", "password")),
		Enabled:  util.IsEnvPresent("AUTH_ENABLED"),
	}
	log.Infof(authConfig.String())

	if util.IsEnvPresent("LIBRATO_TOKEN") {
		go librato.Librato(metrics.DefaultRegistry,
			10e9,                          // interval
			os.Getenv("LIBRATO_EMAIL"),    // account owner email address
			os.Getenv("LIBRATO_TOKEN"),    // Librato API token
			os.Getenv("LIBRATO_HOSTNAME"), // source
			[]float64{0.95},               // percentiles to send
			time.Millisecond,              // time unit
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

	port := util.GetEnvWithDefault("PORT", "8080")
	serverAt := ":" + port
	n.Run(serverAt)
}
