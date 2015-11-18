package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/ello-go/streams/api"
	"github.com/ello/ello-go/streams/model"
	"github.com/julienschmidt/httprouter"
	"github.com/m4rw3r/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var StreamID uuid.UUID

type mockStreamService struct {
	internal []model.StreamItem
}

func (s mockStreamService) AddContent(items []model.StreamItem) error {
	s.internal = append(s.internal, items...)
	return nil
}

func (s mockStreamService) LoadContent(query model.StreamQuery) ([]model.StreamItem, error) {
	return s.internal, nil
}

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

	StreamID, _ := uuid.V4()

	streamService := mockStreamService{
		internal: generateFakeResponse(StreamID),
	}
	streamController := api.NewStreamController(streamService)

	streamController.Register(router)
})

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

func generateFakeResponse(streamID uuid.UUID) []model.StreamItem {
	//fake data
	uuid1, _ := uuid.V4()
	uuid2, _ := uuid.V4()

	return []model.StreamItem{
		{
			ID:        uuid1,
			Timestamp: time.Now(),
			Type:      model.TypePost,
			StreamID:  streamID,
		},
		{
			ID:        uuid2,
			Timestamp: time.Now(),
			Type:      model.TypeRepost,
			StreamID:  streamID,
		},
	}
}

func logResponse(r *httptest.ResponseRecorder) {
	log.WithFields(log.Fields{
		"status":  r.Code,
		"headers": r.HeaderMap,
		"body":    r.Body.String(),
	}).Debug("Got Response")
}
