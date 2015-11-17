package api_test

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ello/ello-go/streams/api"
	"github.com/m4rw3r/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StreamController", func() {
	var id uuid.UUID

	BeforeEach(func() {
		id, _ = uuid.V4()
	})

	Context("when adding content via PUT /streams", func() {

		It("should return a status 201 when passed a correct body", func() {
			item1ID, _ := uuid.V4()
			item2ID, _ := uuid.V4()
			items := []api.StreamItem{{
				StreamID:  id,
				Timestamp: time.Now(),
				Type:      0,
				ID:        item1ID,
			}, {
				StreamID:  id,
				Timestamp: time.Now(),
				Type:      1,
				ID:        item2ID,
			}}
			itemsJSON, _ := json.Marshal(items)
			Request("PUT", "/streams", string(itemsJSON))
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))
			//TODO Verify it tries to add to the StreamService
		})

		It("should return a status 201 when passed a correct body string", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				},
				{
					"id":"c8f17401-62d0-444c-a5d6-639b01f6070f",
					"ts":"2015-11-16T11:59:29.313068877-07:00",
					"type":1,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		It("should return a status 422 when passed an invalid uuid", func() {
			jsonStr := `[
				{
					"id":"ABC",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid date (non ISO8601)", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid type", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":a,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when validation error is in later element", func() {
			jsonStr := `[
				{
					"id":"b8623503-fa3b-4559-9d45-0571a76a98b3",
					"ts":"2015-11-16T11:59:29.313068869-07:00",
					"type":0,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				},
				{
					"id":"c8f17401-62d0-444c-a5d6-639b01f6070f",
					"ts":"2015-11-16T11:59:29.313068877-07:00",
					"type":a,
					"stream_id":"3b1ded01-99ed-4326-9d0b-20127104a2cb"
				}
			]`

			Request("PUT", "/streams", jsonStr)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid body/query", func() {
			Request("PUT", "/streams", "hi")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
	})
	Context("when retrieving a stream via /stream/:id", func() {

		It("should return a status 201 when accessed with a valid ID", func() {
			Request("GET", "/stream/"+id.String(), "")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 422 when passed an invalid id", func() {
			Request("GET", "/stream/"+"abc123", "")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
	})
	Context("when retrieving streams via /streams/coalesce", func() {

		It("should return a status 200 with a valid query string", func() {
			q := api.StreamQuery{
				Streams: []uuid.UUID{id},
			}
			json, _ := json.Marshal(q)
			Request("POST", "/streams/coalesce", string(json))
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 200 with a valid query string", func() {
			q := `{"streams":["10e30ca7-b64d-4510-aaff-775fad0f62ed","6da0fb88-f8f5-40d3-a42c-97147a41011d"]}`
			Request("POST", "/streams/coalesce", q)
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 422 with an invalid uuid", func() {
			q := `{"streams":["10e30ca7-b64d-4510-aaff-775fad0f62ed","abc123"]}`
			Request("POST", "/streams/coalesce", q)
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		It("should return a status 422 when passed an invalid query", func() {
			Request("POST", "/streams/coalesce", "")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})
	})
})
