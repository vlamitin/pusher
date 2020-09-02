package messages

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type MessageSender interface {
	SendMessage(title string, message string) error
}

type CountsReceiver interface {
	GetCountsByTimeFrame(from time.Time, to time.Time) ([]StatusCount, error)
}

type MessageSenderController struct {
	MessageSender MessageSender
}

func NewMessageSenderController(ms MessageSender) *MessageSenderController {
	return &MessageSenderController{MessageSender: ms}
}

type CountsReceiverController struct {
	CountsReceiver CountsReceiver
}

func NewCountsReceiverController(cr CountsReceiver) *CountsReceiverController {
	return &CountsReceiverController{CountsReceiver: cr}
}

type postMessageRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (msc MessageSenderController) PostMessage(w http.ResponseWriter, req *http.Request) {
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

	err = msc.MessageSender.SendMessage(pmr.Title, pmr.Message)
	if err != nil {
		log.Warn(fmt.Sprintf("Error when sending message: %v", err))
		http.Error(w, fmt.Sprintf("Error when sending message"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully sent message")
}

func (crc CountsReceiverController) GetStatusCounts(w http.ResponseWriter, req *http.Request) {
	queries := req.URL.Query()
	from, errFrom := timeFromTimeStampStr(queries.Get("from"))
	to, errTo := timeFromTimeStampStr(queries.Get("to"))
	if errFrom != nil || errTo != nil {
		log.Info(fmt.Sprintf("Error when parsing query values from from client: %v, %v", errFrom, errTo))
		http.Error(w, fmt.Sprintf("No `from` or `to` query params receieved"), http.StatusBadRequest)
		return
	}

	counts, err := crc.CountsReceiver.GetCountsByTimeFrame(from, to)
	if err != nil {
		log.Warn(fmt.Sprintf("Error when querying db for counts: %v", err))
		http.Error(w, fmt.Sprintf("Failed to query counts"), http.StatusInternalServerError)
		return
	}

	jData, err := json.Marshal(counts)
	if err != nil {
		log.Warn(fmt.Sprintf("Error when encoding counst: %v", err))
		http.Error(w, fmt.Sprintf("Failed to create json"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func timeFromTimeStampStr(stamp string) (time.Time, error) {
	integer, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Unix(integer/1000, (integer%1000)*1000000), nil
}
