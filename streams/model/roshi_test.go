package model_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/ello-go/streams/model"
	"github.com/m4rw3r/uuid"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.WarnLevel)

	retCode := m.Run()

	os.Exit(retCode)
}

func TestJsonMarshal(t *testing.T) {
	id, _ := uuid.V4()
	item := model.StreamItem{
		StreamID:  id,
		Timestamp: time.Now(),
		Type:      0,
		ID:        id,
	}

	output, _ := json.Marshal(item)
	var fromJSON model.StreamItem
	err := json.Unmarshal(output, &fromJSON)

	log.WithFields(log.Fields{
		"Source":   item,
		"JSON":     string(output),
		"fromJSON": fromJSON,
		"ERROR":    err,
	}).Debug("StreamItem Example")

	if item != fromJSON {
		t.Errorf("Source doesn't match the marshal/unmarshaled value")
	}

	output2, _ := json.Marshal(model.RoshiStreamItem(item))
	var fromJSON2 model.RoshiStreamItem
	err = json.Unmarshal(output2, &fromJSON2)

	log.WithFields(log.Fields{
		"Source":   item,
		"JSON":     string(output2),
		"fromJSON": fromJSON2,
		"ERROR":    err,
	}).Debug("RoshiStreamItem Example")

	if model.RoshiStreamItem(item) != fromJSON2 {
		t.Errorf("Source doesn't match the marshal/unmarshaled value with RoshiStreamItem")
	}

	items := []model.StreamItem{
		item,
		item,
	}
	rItems, _ := model.MarshalRoshi(items)
	output3, _ := json.Marshal(rItems)
	var fromJSON3 []model.RoshiStreamItem
	err = json.Unmarshal(output3, &fromJSON3)

	log.WithFields(log.Fields{
		"Source":   rItems,
		"JSON":     string(output3),
		"fromJSON": fromJSON3,
		"ERROR":    err,
	}).Debug("RoshiStreamItem Example")

	if !reflect.DeepEqual(rItems, fromJSON3) {
		t.Errorf("Source doesn't match the marshal/unmarshaled value with []RoshiStreamItem")
	}
}
