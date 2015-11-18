package model

import (
	"time"

	"github.com/m4rw3r/uuid"
)

// StreamItemType does stuff
type StreamItemType int

const (
	//TypePost is a type of stream item which is a direct post
	TypePost StreamItemType = iota
	//TypeRepost is a type of stream item which represents a repost
	TypeRepost
)

//StreamItem represents a single item on a stream
type StreamItem struct {
	ID        uuid.UUID      `json:"id"`
	Timestamp time.Time      `json:"ts"`
	Type      StreamItemType `json:"type"`
	StreamID  uuid.UUID      `json:"stream_id"`
}

//StreamQuery represents a query for multiple streams
type StreamQuery struct {
	Streams []uuid.UUID `json:"streams"`
}
