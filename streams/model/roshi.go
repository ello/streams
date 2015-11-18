package model

import (
	"encoding/json"
	"time"

	"github.com/m4rw3r/uuid"
)

type roshiMember struct {
	ID   uuid.UUID      `json:"content_id"`
	Type StreamItemType `json:"type"`
}

type roshiInsert struct {
	Key    []byte `json:"key"`
	Score  int64  `json:"score"`
	Member []byte `json:"member"`
}

//RoshiStreamItem shadows StreamItem to allow us to export the json Roshi expects
type RoshiStreamItem StreamItem

// MarshalJSON converts from a RoshiStreamItem to the expected json for Roshi
func (item RoshiStreamItem) MarshalJSON() ([]byte, error) {
	member, _ := json.Marshal(&roshiMember{
		ID:   item.ID,
		Type: item.Type,
	})
	return json.Marshal(&roshiInsert{
		Key:    []byte(item.StreamID.String()),
		Score:  item.Timestamp.UnixNano(),
		Member: []byte(member),
	})
}

//UnmarshalJSON correct converts a roshi json blob back to RoshiStreamItem
func (item *RoshiStreamItem) UnmarshalJSON(data []byte) error {
	var jsonItem roshiInsert
	err := json.Unmarshal(data, &jsonItem)
	if err == nil {
		//unpack the streamID back to a UUID
		streamID, innerErr := uuid.FromString(string(jsonItem.Key))
		if innerErr != nil {
			return innerErr
		}

		//unpack the body of the record for the id and type
		var member roshiMember
		innerErr = json.Unmarshal(jsonItem.Member, &member)
		if innerErr != nil {
			return innerErr
		}

		//set the values
		item.StreamID = streamID
		item.Timestamp = time.Unix(0, jsonItem.Score)
		item.Type = member.Type
		item.ID = member.ID

	}
	return err
}

//MarshalRoshi converts a slice of StreamItems into a slice of RoshiStreamItems
func MarshalRoshi(items []StreamItem) ([]RoshiStreamItem, error) {
	rItems := make([]RoshiStreamItem, len(items))
	for i := 0; i < len(items); i++ {
		rItems[i] = RoshiStreamItem(items[i])
	}
	return rItems, nil
}

//UnmarshalRoshi converts a slice of RoshiStreamItems to a slice of StreamItems
func UnmarshalRoshi(rItems []RoshiStreamItem) ([]StreamItem, error) {
	items := make([]StreamItem, len(rItems))
	for i := 0; i < len(rItems); i++ {
		items[i] = StreamItem(rItems[i])
	}
	return items, nil
}
