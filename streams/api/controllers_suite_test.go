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
	internal       []model.StreamItem
	lastItemsOnAdd []model.StreamItem
	lastLimit      int
	lastFromSlug   string
}

func (s mockStreamService) Add(items []model.StreamItem) error {
	s.lastItemsOnAdd = items
	s.internal = append(s.internal, items...)
	return nil
}

func (s mockStreamService) Load(query model.StreamQuery, limit int, fromSlug string) (*model.StreamQueryResponse, error) {
	s.lastLimit = limit
	s.lastFromSlug = fromSlug
	return &model.StreamQueryResponse{Items: s.internal}, nil
}

var (
	router        *httprouter.Router
	response      *httptest.ResponseRecorder
	streamService mockStreamService
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

	streamService = mockStreamService{
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

func checkStreamItems(c model.StreamItem, c1 model.StreamItem) {
	Expect(c).NotTo(BeNil())
	Expect(c1).NotTo(BeNil())
	Expect(c1.StreamID).To(Equal(c.StreamID))
	Expect(c1.ID).To(Equal(c.ID))
	Expect(c1.Type).To(Equal(c.Type))
	Expect(c1.Timestamp).To(BeTemporally("~", c.Timestamp, time.Millisecond))
}

func checkAll(c []model.StreamItem, c1 []model.StreamItem) {
	for i := 0; i < len(c); i++ {
		checkStreamItems(c[i], c1[i])
	}
}
