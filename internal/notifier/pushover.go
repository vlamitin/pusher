package notifier

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type PushoverConnector struct {
	userToken string
	appToken  string
}

func NewPushoverConnector(userToken string, appToken string) *PushoverConnector {
	return &PushoverConnector{
		userToken: userToken,
		appToken:  appToken,
	}
}

type notifyResult struct {
	Status int `json:"status"`
}

func (pc *PushoverConnector) Notify(title string, message string) (status int, err error) {
	data := url.Values{}
	data.Set("token", pc.appToken)
	data.Set("user", pc.userToken)
	data.Set("title", title)
	data.Set("message", message)

	log.Infof("Sending title %s and message %s", title, message)

	req, err := http.NewRequest(http.MethodPost, "https://api.pushover.net/1/messages.json", strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return -1, fmt.Errorf("error preparing pushover request: %v\n", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("error when performing pushover request: %v\n", err)
	}
	defer resp.Body.Close()

	var jsonResponse notifyResult
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&jsonResponse)
	if err != nil {
		return -1, fmt.Errorf("error when decoding pushover response: %v\n", err)
	}

	log.Infof("Successfully sent with status %d", jsonResponse.Status)

	return jsonResponse.Status, nil
}
