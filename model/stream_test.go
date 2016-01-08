package model_test

import (
	"testing"
	"time"

	"github.com/ello/streams/model"
)

func CheckStreamItems(c model.StreamItem, c1 model.StreamItem, t *testing.T) {
	if c1.StreamID != c.StreamID {
		t.Error("StreamIDs should match")
	}
	if c1.ID != c.ID {
		t.Error("IDs should match")
	}
	if c1.Type != c.Type {
		t.Error("Type should match")
	}
	if !c1.Timestamp.Round(time.Millisecond).Equal(c.Timestamp.Round(time.Millisecond)) {
		t.Error("Timestamps should be within a millisecond")
	}
}

func CheckAll(c []model.StreamItem, c1 []model.StreamItem, t *testing.T) {
	for i := 0; i < len(c); i++ {
		CheckStreamItems(c[i], c1[i], t)
	}
}
