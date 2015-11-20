package service

import "github.com/ello/ello-go/streams/model"

// StreamService is the interface to the underlying stream storage system.
type StreamService interface {

	//Add will add the content items to the embedded stream id
	Add(items []model.StreamItem) error

	//Load will pull a coalesced view of the streams in the query
	Load(query model.StreamQuery, limit int, offset int) ([]model.StreamItem, error)
}
