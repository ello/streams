package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/ello/ello-go/streams/model"
)

//NewRoshiStreamService takes a url for the roshi server and returns the service
func NewRoshiStreamService(urlString string, timeoutSeconds time.Duration) (StreamService, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return roshiStreamService{
		url:     u,
		timeout: timeoutSeconds * time.Second,
	}, nil
}

type roshiStreamService struct {
	url     *url.URL
	timeout time.Duration
}

func (s roshiStreamService) Add(items []model.StreamItem) error {
	rItems, err := model.ToRoshiStreamItem(items)
	if err != nil {
		log.Error(err)
		return err
	}

	requestBody, err := json.Marshal(rItems)
	if err != nil {
		log.Error(err)
		return err
	}

	uri := s.url.String()

	log.WithFields(log.Fields{
		"Body": string(requestBody),
		"URL":  uri,
	}).Debug("Preparing to make request")

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(requestBody))
	if log.GetLevel() >= log.DebugLevel {
		debug(httputil.DumpRequestOut(req, true))
	}

	if err != nil {
		log.Error(err)
		return err
	}
	client := &http.Client{
		Timeout: s.timeout,
	}
	log.WithFields(log.Fields{
		"client": client,
		"req":    req,
	}).Debug("About to execute")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		debug(httputil.DumpResponse(resp, true))
		return errors.New("Request Failed with status: " + string(resp.StatusCode))
	}

	return nil
}

func (s roshiStreamService) Load(query model.StreamQuery, limit int, cursor string) (*model.StreamQueryResponse, error) {
	requestBody, err := json.Marshal(model.RoshiQuery(query))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	uri := fmt.Sprintf("%v?coalesce=true&limit=%d", s.url.String(), limit)
	if len(cursor) != 0 {
		// TODO Should probably validate the slug is valid here and return an error if not
		uri = fmt.Sprintf("%v&start=%v", uri, cursor)
	}

	log.WithFields(log.Fields{
		"Body": string(requestBody),
		"URL":  uri,
	}).Debug("Preparing to make request")

	req, err := http.NewRequest("GET", uri, bytes.NewBuffer(requestBody))
	if log.GetLevel() >= log.DebugLevel {
		debug(httputil.DumpRequestOut(req, true))
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}
	client := &http.Client{
		Timeout: s.timeout,
	}
	log.WithFields(log.Fields{
		"client": client,
		"req":    req,
	}).Debug("About to execute")

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 || err != nil {
		debug(httputil.DumpResponse(resp, true))
		return nil, errors.New("Request Failed with status: " + string(resp.StatusCode))
	}

	var result model.RoshiResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Debugf("Data: %v", string(data))
		log.Errorf("Error unmarshalling result: %v", err)
		return nil, err
	}

	log.WithFields(log.Fields{
		"Status":   resp.StatusCode,
		"Duration": result.Duration,
		"Records":  result.Items,
		"Raw":      string(data),
	}).Debug("Execution complete")

	items, err := model.ToStreamItem(result.Items)

	return &model.StreamQueryResponse{
		Items:  items,
		Cursor: generateCursor(result.Items),
	}, err
}

func generateCursor(items []model.RoshiStreamItem) string {
	if len(items) == 0 {
		return ""
	}
	oldest := items[len(items)-1]

	ts := oldest.Timestamp
	tsBits := math.Float64bits(float64(ts.UnixNano()))
	member, _ := model.MemberJSON(oldest)
	encodedMember := base64.StdEncoding.EncodeToString(member)
	cursor := fmt.Sprintf("%dA%s", tsBits, encodedMember)

	log.WithFields(log.Fields{
		"Time":           ts,
		"Time in Bits":   tsBits,
		"Member":         string(member),
		"Encoded Member": encodedMember,
		"Cursor":         cursor,
	}).Debug("Generated Cursor")

	return cursor
}

func debug(data []byte, err error) {
	log.WithFields(log.Fields{
		"Req/Res": fmt.Sprintf("\n%s\n\n", data),
		"Error":   err,
	}).Debug("Debugging")
}
