package service

import "github.com/ello/ello-go/streams/model"

// StreamService is the interface to the underlying stream storage system.
type StreamService interface {
	AddContent(items []model.StreamItem) error
	LoadContent(query model.StreamQuery) ([]model.StreamItem, error)
}
