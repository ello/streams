package model_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ello/streams/model"
	"github.com/m4rw3r/uuid"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)

	retCode := m.Run()

	os.Exit(retCode)
}

func TestJsonMarshal(t *testing.T) {
	id, _ := uuid.V4()
	item := model.StreamItem{
		StreamID: id.String(),
		// NOTE:  Converting between int64/float64 at the nanosecond level can
		//				create some tiny drift. In practice, this is fine.  Rounding to
		//				the second level avoids issues with testing.
		Timestamp: time.Now().Round(time.Second),
		Type:      model.TypePost,
		ID:        id.String(),
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

	if err != nil {
		t.Errorf("Error converting to/from json")
	}

	CheckStreamItems(item, fromJSON, t)

	output2, _ := json.Marshal(model.RoshiStreamItem(item))
	var fromJSON2 model.RoshiStreamItem
	err = json.Unmarshal(output2, &fromJSON2)

	log.WithFields(log.Fields{
		"Source":   item,
		"JSON":     string(output2),
		"fromJSON": fromJSON2,
		"ERROR":    err,
	}).Debug("RoshiStreamItem Example")

	if err != nil {
		t.Errorf("Error converting to/from json")
	}

	checkRoshiItems(model.RoshiStreamItem(item), fromJSON2, t)

	items := []model.StreamItem{
		item,
		item,
	}
	rItems, _ := model.ToRoshiStreamItem(items)
	output3, _ := json.Marshal(rItems)
	var fromJSON3 []model.RoshiStreamItem
	err = json.Unmarshal(output3, &fromJSON3)

	log.WithFields(log.Fields{
		"Source":   rItems,
		"JSON":     string(output3),
		"fromJSON": fromJSON3,
		"ERROR":    err,
	}).Debug("RoshiStreamItem Example")

	checkAllRoshi(rItems, fromJSON3, t)
}

func checkRoshiItems(c model.RoshiStreamItem, c1 model.RoshiStreamItem, t *testing.T) {
	c.StreamID = string(c.StreamID)
	CheckStreamItems(model.StreamItem(c), model.StreamItem(c1), t)
}

func checkAllRoshi(c []model.RoshiStreamItem, c1 []model.RoshiStreamItem, t *testing.T) {
	for i := 0; i < len(c); i++ {
		checkRoshiItems(c[i], c1[i], t)
	}
}
