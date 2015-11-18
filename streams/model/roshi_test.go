package model_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ello/ello-go/streams/model"
	"github.com/m4rw3r/uuid"
)

func TestJsonMarshal(t *testing.T) {

	id, _ := uuid.V4()
	item := model.StreamItem{
		StreamID:  id,
		Timestamp: time.Now(),
		Type:      0,
		ID:        id,
	}

	output, _ := json.Marshal(item)

	output2, _ := json.Marshal(model.RoshiStreamItem(item))

	fmt.Println(string(output))
	fmt.Println(string(output2))

	items := []model.StreamItem{
		item,
		item,
	}

	rItems, _ := model.MarshalRoshi(items)

	fmt.Println(rItems)

	output3, _ := json.Marshal(rItems)
	fmt.Println(string(output3))

}
