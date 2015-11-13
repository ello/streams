package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/ello-go/streams/controllers"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

//
// type mockUserService struct {
// 	lastOffset int
// 	lastLimit  int
// }
//
// func (s *mockUserService) FindByUsername(username string) (service.User, error) {
// 	if username == "rtyer" {
// 		return service.User{Username: "rtyer", Email: "rtyer@ello.co", ID: 12345}, nil
// 	}
// 	return service.User{}, errors.New("Username not found")
//
// }
// func (s *mockUserService) FindUsers(limit int, offset int) ([]service.User, error) {
// 	s.lastLimit = limit
// 	s.lastOffset = offset
// 	return []service.User{
// 		service.User{Username: "rtyer", Email: "rtyer@ello.co", ID: 12345},
// 	}, nil
// }

var (
	router   *httprouter.Router
	response *httptest.ResponseRecorder
)

func Request(method string, route string, body string) {
	r, err := http.NewRequest(method, route, strings.NewReader(body))
	response = httptest.NewRecorder()

	log.WithFields(log.Fields{
		"url":     r.URL,
		"method":  r.Method,
		"headers": r.Header,
		"body":    body,
		"errors":  err,
	}).Debug("About to issue request")

	router.ServeHTTP(response, r)
}

var _ = BeforeSuite(func() {
	log.SetLevel(log.DebugLevel)

	router = httprouter.New()

	streamController := controllers.NewStreamController()

	streamController.Register(router)
})

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}
