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
)

var commit string
var startTime = time.Now()
var help bool
var helpMessage = `ELLO STREAM API
--------------------------
Set ENV Variables to configure:

ELLO_API_PORT for the port to run this service on.  Default is 8080
ELLO_ROSHI_HOST for the location of the roshi instance.  Default is http://localhost:6302
ELLO_AUTH_ENABLED any value will enable basic auth.  Default is disabled.
ELLO_AUTH_USERNAME for the auth username.  Default is 'ello'.
ELLO_AUTH_PASSWORD for the auth password.  Default is 'password'.
ELLO_LOG_LEVEL for the log level.  Valid levels are "debug", "info", "warn", "error"
`

func main() {

	flag.BoolVar(&help, "h", false, "help?")
	flag.Parse()

	if help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}
	level := util.GetEnvWithDefault("ELLO_LOG_LEVEL", "http://localhost:6302")

	logLevel := log.WarnLevel
	switch level {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "error":
		logLevel = log.ErrorLevel
	}
	log.SetLevel(logLevel)
	fmt.Printf("Using log level [%v]\n", logLevel)

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
