package service

import "github.com/ello/ello-go/streams/model"

// StreamService does shit
type StreamService interface {
	AddContent(items []model.StreamItem) error
}
