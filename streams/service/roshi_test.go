package service_test

import (
	"time"

	"github.com/ello/ello-go/streams/model"
	"github.com/ello/ello-go/streams/service"
	"github.com/m4rw3r/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Roshi Channel Service", func() {
	var _ = Describe("Instantiation", func() {

		It("sanity?", func() {
			s, err := service.NewRoshiStreamService("http://localhost:6302")
			Expect(err).To(BeNil())
			Expect(s).NotTo(BeNil())
		})

	})
	var s service.StreamService
	BeforeEach(func() {
		s, _ = service.NewRoshiStreamService("http://localhost:6302")
	})

	AfterEach(func() {

	})

	It("is sane", func() {
		Expect(true).Should(BeTrue())
	})

	Context(".AddToChannel", func() {
		It("will add a single content item", func() {
			chanID, _ := uuid.V4()
			contentID, _ := uuid.V4()

			content := model.StreamItem{
				ID:        chanID,
				Timestamp: time.Now(),
				Type:      model.TypePost,
				StreamID:  contentID,
			}
			items := []model.StreamItem{
				content,
			}
			Expect(s).NotTo(BeNil())
			err := s.AddContent(items)
			Expect(err).To(BeNil())
		})

		// Context(".LoadChannel", func() {
		// 	It("Load content previously added to the channel", func() {
		// 		chanID, _ := uuid.V4()
		// 		contentID, _ := uuid.V4()
		// 		content := data.Content{
		// 			ID:        contentID,
		// 			CreatedAt: time.Now(),
		// 		}
		//
		// 		s.AddToChannel(chanID, content)
		//
		// 		contentIDs, _ := s.LoadChannel(chanID)
		// 		Expect(contentIDs).NotTo(BeEmpty())
		// 	})
		// })
	})
})
