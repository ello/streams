package controllers_test

import (
	"net/http"

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

			Request("PUT", "/streams", "hi")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		It("should attempt to add the content item to the streamservice", func() {
			//todo
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

		It("should return a status 201 when accessed with a valid ID", func() {
			Request("POST", "/streams/coalesce", "")
			logResponse(response)

			Expect(response.Code).To(Equal(http.StatusOK))
		})

		It("should return a status 422 when passed an invalid query", func() {
			Request("POST", "/streams/coalesce", "")
			logResponse(response)

			Expect(response.Code).To(Equal(422))
		})

		// 	It("should return a status 200 with no args", func() {
		// 		Request("GET", "/users")
		// 		var data []service.User
		// 		_ = json.Unmarshal(response.Body.Bytes(), &data)
		//
		// 		Expect(response.Code).To(Equal(http.StatusOK))
		// 		Expect(data[0].Username).To(Equal("rtyer"))
		// 	})
		//
		// 	It("should use the passed limit/offset", func() {
		// 		Request("GET", "/users?limit=5&offset=13")
		//
		// 		Expect(response.Code).To(Equal(http.StatusOK))
		// 		Expect(userService.lastLimit).To(Equal(5))
		// 		Expect(userService.lastOffset).To(Equal(13))
		// 	})
		//
		// 	It("should correctly validate the limit", func() {
		// 		Request("GET", "/users?limit=a")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotAcceptable))
		// 	})
		//
		// 	It("should correctly validate the offset", func() {
		// 		Request("GET", "/users?offset=a")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotAcceptable))
		// 	})
		// })
		// Context("when calling /user/<username>", func() {
		//
		// 	It("should return a status 200 with a user that is present", func() {
		// 		Request("GET", "/users/rtyer")
		// 		var user service.User
		// 		_ = json.Unmarshal(response.Body.Bytes(), &user)
		//
		// 		Expect(response.Code).To(Equal(http.StatusOK))
		// 		Expect(user.Username).To(Equal("rtyer"))
		// 	})
		//
		// 	It("should return a status 404 with a non existent user", func() {
		// 		Request("GET", "/users/asdf")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotFound))
		// 	})
		//
		// 	It("should return a status 406 if the username is invalid", func() {
		// 		Request("GET", "/users/^&*$")
		//
		// 		Expect(response.Code).To(Equal(http.StatusNotAcceptable))
		// 	})
		//
		// 	It("should accept UTF-8 characters for the username", func() {
		// 		Request("GET", "/users/ßåœ")
		// 		Expect(response.Code).NotTo(Equal(http.StatusNotAcceptable))
		// 	})
	})
})
