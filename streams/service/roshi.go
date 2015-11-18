package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/Sirupsen/logrus"

	"github.com/ello/ello-go/streams/model"
)

//NewRoshiStreamService takes a url for the roshi server and returns the service
func NewRoshiStreamService(urlString string) (StreamService, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return roshiStreamService{
		url: u,
	}, nil
}

type roshiStreamService struct {
	url *url.URL
}

func (s roshiStreamService) AddContent(items []model.StreamItem) error {
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
	client := &http.Client{}
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

func (s roshiStreamService) LoadContent(query model.StreamQuery) ([]model.StreamItem, error) {
	requestBody, err := json.Marshal(model.RoshiQuery(query))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	uri := s.url.String() + "?coalesce=true"

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
	client := &http.Client{}
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
		log.Error(err)
		return nil, err
	}

	log.WithFields(log.Fields{
		"Status":   resp.StatusCode,
		"Duration": result.Duration,
		"Records":  result.Items,
		"Raw":      string(data),
	}).Debug("Execution complete")

	return model.ToStreamItem(result.Items)
}

func debug(data []byte, err error) {
	log.WithFields(log.Fields{
		"Req/Res": fmt.Sprintf("\n%s\n\n", data),
		"Error":   err,
	}).Debug("Debugging")
}
