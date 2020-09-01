package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type TesterConfig struct {
	serverHost string
	serverPort int
}

func main() {
	config := getTesterConfig()
	TestSendNotificationScenario(config.serverHost, config.serverPort)
}

func getTesterConfig() TesterConfig {
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	return TesterConfig{
		serverHost: os.Getenv("SERVER_HOST"),
		serverPort: serverPort,
	}
}

func TestSendNotificationScenario(serverHost string, serverPort int) {
	var jsonStr = []byte(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, "title", "message"))
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/messages", serverHost, serverPort), bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	expected := "Successfully sent message"

	if body := string(bodyBytes); body != expected {
		log.Fatal(fmt.Errorf("handler returned unexpected body: got `%v` want `%v`", body, expected))
	}
}
