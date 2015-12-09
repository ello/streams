package model

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/OneOfOne/xxhash"
)

type roshiBody struct {
	ID       string         `json:"content_id"`
	StreamID string         `json:"stream_id"`
	Type     StreamItemType `json:"type"`
}

type roshiItem struct {
	Key    []byte  `json:"key"`
	Score  float64 `json:"score"`
	Member []byte  `json:"member"`
}

//RoshiResponse represents the response from a Query
type RoshiResponse struct {
	Duration string            `json:"duration"`
	Items    []RoshiStreamItem `json:"records"`
}

//RoshiStreamItem shadows StreamItem to allow us to export the json Roshi expects
type RoshiStreamItem StreamItem

//RoshiQuery shadows a StreamItem to allow us to export the json Roshi expects
type RoshiQuery StreamQuery

// MarshalJSON converts from a RoshiStreamItem to the expected json for Roshi
func (item RoshiStreamItem) MarshalJSON() ([]byte, error) {
	member, _ := MemberJSON(item)
	h := xxhash.New64()
	io.Copy(h, strings.NewReader(item.StreamID))
	return json.Marshal(&roshiItem{
		Key:    []byte(h.Sum(nil)),
		Score:  float64(item.Timestamp.UnixNano()),
		Member: []byte(member),
	})
}

//MemberJSON Returns the byte array of the json for a given stream item in roshi member form
func MemberJSON(item RoshiStreamItem) ([]byte, error) {
	return json.Marshal(&roshiBody{
		ID:       item.ID,
		Type:     item.Type,
		StreamID: item.StreamID,
	})
}

//UnmarshalJSON correct converts a roshi json blob back to RoshiStreamItem
func (item *RoshiStreamItem) UnmarshalJSON(data []byte) error {
	var jsonItem roshiItem
	err := json.Unmarshal(data, &jsonItem)
	if err == nil {
		//unpack the body of the record for the id and type
		var member roshiBody
		innerErr := json.Unmarshal(jsonItem.Member, &member)
		if innerErr != nil {
			return innerErr
		}

		//set the values
		item.StreamID = member.StreamID
		item.Timestamp = time.Unix(0, int64(jsonItem.Score))
		item.Type = member.Type
		item.ID = member.ID

	}
	return err
}

//MarshalJSON takes a roshiquery and creates a list of base64 encoded bytes
func (q RoshiQuery) MarshalJSON() ([]byte, error) {
	ids := make([][]byte, len(q.Streams))
	for i := 0; i < len(q.Streams); i++ {
		h := xxhash.New64()
		io.Copy(h, strings.NewReader(q.Streams[i]))
		ids[i] = h.Sum(nil)
	}
	return json.Marshal(ids)
}

//ToRoshiStreamItem converts a slice of StreamItems into a slice of RoshiStreamItems
func ToRoshiStreamItem(items []StreamItem) ([]RoshiStreamItem, error) {
	rItems := make([]RoshiStreamItem, len(items))
	for i := 0; i < len(items); i++ {
		rItems[i] = RoshiStreamItem(items[i])
	}
	return rItems, nil
}

//ToStreamItem converts a slice of RoshiStreamItems to a slice of StreamItems
func ToStreamItem(rItems []RoshiStreamItem) ([]StreamItem, error) {
	items := make([]StreamItem, len(rItems))
	for i := 0; i < len(rItems); i++ {
		items[i] = StreamItem(rItems[i])
	}
	return items, nil
}
