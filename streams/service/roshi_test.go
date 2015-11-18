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

	Context(".AddToChannel", func() {
		It("will add a single content item", func() {
			chanID, _ := uuid.V4()
			contentID, _ := uuid.V4()

			content := model.StreamItem{
				ID:        contentID,
				Timestamp: time.Now(),
				Type:      model.TypePost,
				StreamID:  chanID,
			}
			items := []model.StreamItem{
				content,
			}
			Expect(s).NotTo(BeNil())
			err := s.AddContent(items)
			Expect(err).To(BeNil())
		})

		Context(".LoadChannel", func() {
			It("Load content previously added to the channel", func() {
				chanID, _ := uuid.V4()
				contentID, _ := uuid.V4()

				content := model.StreamItem{
					ID:        contentID,
					Timestamp: time.Now(),
					Type:      model.TypePost,
					StreamID:  chanID,
				}
				items := []model.StreamItem{
					content,
				}
				err := s.AddContent(items)
				Expect(err).To(BeNil())

				itemID, _ := uuid.V4()
				q := model.StreamQuery{
					Streams: []uuid.UUID{itemID, chanID},
				}

				c, _ := s.LoadContent(q)
				Expect(c).NotTo(BeEmpty())
				Expect(len(c)).To(Equal(1))
				c1 := c[0]
				Expect(c1.StreamID).To(Equal(content.StreamID))
				Expect(c1.ID).To(Equal(content.ID))
				Expect(c1.Type).To(Equal(content.Type))
				Expect(c1.Timestamp).To(BeTemporally("~", content.Timestamp, time.Millisecond))
			})
		})
	})
})
