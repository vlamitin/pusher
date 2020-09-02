package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TesterConfig struct {
	serverHost string
	serverPort int
}

func main() {
	config := getTesterConfig()
	TestSendNotificationAndQueryStatusesScenario(fmt.Sprintf("http://%s:%d", config.serverHost, config.serverPort))
	log.Info("No tests failed")
}

func getTesterConfig() TesterConfig {
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	return TesterConfig{
		serverHost: os.Getenv("SERVER_HOST"),
		serverPort: serverPort,
	}
}

func TestSendNotificationAndQueryStatusesScenario(baseUrl string) {
	beforeSendMessage := time.Now()
	sendMessageRes := sendMessage(baseUrl)
	sendMessageExp := "Successfully sent message"
	if sendMessageRes != sendMessageExp {
		log.Fatal(fmt.Errorf("handler returned unexpected body: got `%v` want `%v`", sendMessageRes, sendMessageExp))
	}

	queryCountsRes := queryCounts(baseUrl, beforeSendMessage, time.Now())
	queryCountsExp := `[{"SendStatus":1,"Count":1}]`
	if queryCountsRes != queryCountsExp {
		log.Fatal(fmt.Errorf("handler returned unexpected body: got `%v` want `%v`", queryCountsRes, queryCountsExp))
	}
}

func sendMessage(baseUrl string) string {
	var jsonStr = []byte(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, "title", "message"))
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/messages", baseUrl), bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return string(bodyBytes)
}

func queryCounts(baseUrl string, from time.Time, to time.Time) string {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/messages?from=%d&to=%d", baseUrl, from.UnixNano()/1000000, to.UnixNano()/1000000), nil)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	return string(bodyBytes)
}
