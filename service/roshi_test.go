package service_test

import (
	"time"

	"github.com/ello/streams/model"
	"github.com/ello/streams/service"
	"github.com/ello/streams/util"
	"github.com/m4rw3r/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Roshi Channel Service", func() {
	var _ = Describe("Instantiation", func() {

		It("sanity?", func() {
			s, err := service.NewRoshiStreamService(util.GetEnvWithDefault("ROSHI_URL", "http://localhost:6302"), (5 * time.Second))
			Expect(err).To(BeNil())
			Expect(s).NotTo(BeNil())
		})

	})
	var s service.StreamService
	BeforeEach(func() {
		s, _ = service.NewRoshiStreamService(util.GetEnvWithDefault("ROSHI_URL", "http://localhost:6302"), (5 * time.Second))
	})

	Context(".Add", func() {
		It("will add a single content item", func() {
			chanID, _ := uuid.V4()
			contentID, _ := uuid.V4()

			content := model.StreamItem{
				ID:        contentID.String(),
				Timestamp: time.Now(),
				Type:      model.TypePost,
				StreamID:  chanID.String(),
			}
			items := []model.StreamItem{
				content,
			}
			Expect(s).NotTo(BeNil())
			err := s.Add(items)
			Expect(err).To(BeNil())
		})

		Context(".Load", func() {
			It("Load content previously added to the channel", func() {
				chanID, _ := uuid.V4()
				contentID, _ := uuid.V4()

				content := model.StreamItem{
					ID:        contentID.String(),
					Timestamp: time.Now(),
					Type:      model.TypePost,
					StreamID:  chanID.String(),
				}
				items := []model.StreamItem{
					content,
				}
				err := s.Add(items)
				Expect(err).To(BeNil())

				fakeChanID, _ := uuid.V4()
				q := model.StreamQuery{
					Streams: []string{fakeChanID.String(), chanID.String()},
				}

				resp, _ := s.Load(q, 10, "")
				c := resp.Items
				Expect(c).NotTo(BeEmpty())
				Expect(len(c)).To(Equal(1))
				c1 := c[0]

				Expect(c1.StreamID).To(Equal(content.StreamID))
				Expect(c1.ID).To(Equal(content.ID))
				Expect(c1.Type).To(Equal(content.Type))
				Expect(c1.Timestamp).To(BeTemporally("~", content.Timestamp, time.Millisecond))
			})
		})

		Context(".Remove", func() {
			It("Remove content previously added to the channel", func() {
				chanID, _ := uuid.V4()
				contentID, _ := uuid.V4()

				content := model.StreamItem{
					ID:        contentID.String(),
					Timestamp: time.Now(),
					Type:      model.TypePost,
					StreamID:  chanID.String(),
				}
				items := []model.StreamItem{
					content,
				}
				err := s.Add(items)
				Expect(err).To(BeNil())

				fakeChanID, _ := uuid.V4()
				q := model.StreamQuery{
					Streams: []string{fakeChanID.String(), chanID.String()},
				}

				resp, _ := s.Load(q, 10, "")
				c := resp.Items
				Expect(c).NotTo(BeEmpty())
				Expect(len(c)).To(Equal(1))

				rm_err := s.Remove(items)
				Expect(rm_err).To(BeNil())

				resp2, _ := s.Load(q, 10, "")
				Expect(resp2.Items).To(BeEmpty())
			})
		})
	})
})
