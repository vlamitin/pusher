package messages

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MessageSender interface {
	SendMessage(title string, message string) error
}

type Controller struct {
	MessageSender MessageSender
}

func NewController(ms MessageSender) *Controller {
	return &Controller{MessageSender: ms}
}

type postMessageRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (mc Controller) PostMessage(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, 1048576)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var pmr postMessageRequest
	err := dec.Decode(&pmr)
	if err != nil {
		log.Info(fmt.Sprintf("Error when decoding json from client: %v", err))
		http.Error(w, fmt.Sprintf("Invalid json received"), http.StatusBadRequest)
		return
	}

	err = mc.MessageSender.SendMessage(pmr.Title, pmr.Message)
	if err != nil {
		log.Warn(fmt.Sprintf("Error when sending message: %v", err))
		http.Error(w, fmt.Sprintf("Error when sending message"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully sent message")
}
