package task

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
)

var (
	host          = "127.0.0.1"
	RegisterRoute = "register"
	FeedRoute     = "feedback"
)

func register(body []TaskBase) error {
	return send(RegisterRoute, body)
}

func feedBack(v taskFeedBack) error {
	return send(FeedRoute, v)
}

func send(route string, v interface{}) error {
	jsonBody, err := json.Marshal(v)
	if err != nil {
		return err
	}
	sendBody := bytes.NewReader(jsonBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", path.Join(host, route), sendBody)
	if err != nil {
		return err
	}
	// req.Header.Add("Authorization", "APPCODE "+appcode)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}
