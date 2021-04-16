package task

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	host            = "http://127.0.0.1:8080/task/back/"
	RegisterRoute   = "register"
	FeedRoute       = "feedback"
	StartRoute      = "start"
	EndFeedRoute    = "end"
	UpdateFeedRoute = "update"
)

type MesHandleBack interface{}

func register(body argRegister) error {
	return sendPost(RegisterRoute, body)
}

func feedBack(v taskFeedBack) error {
	return sendPost(FeedRoute, v)
}

// func start(body argTask) error {
// 	return sendPost(StartRoute, body)
// }

func end(body argTask) error {
	return sendPost(EndFeedRoute, body)
}

// func update(body argTaskUpdate) error {
// 	return sendPost(UpdateFeedRoute, body)
// }

func sendPost(route string, v MesHandleBack) error {
	jsonBody, err := json.Marshal(v)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", host+route, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	// req.Header.Add("Authorization", "APPCODE "+appcode)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
