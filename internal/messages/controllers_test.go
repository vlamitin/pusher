package messages

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostMessageSuccess(t *testing.T) {
	requestBodyStr := []byte(`{"title": "title", "message": "message"}`)
	req, err := http.NewRequest("POST", "/messages", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	controller := NewMessageSenderController(&fakeMessageSenderSuccess{})
	controller.PostMessage(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got `%v` want `%v`",
			status, http.StatusOK)
	}

	expected := "Successfully sent message"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got `%v` want `%v`",
			rr.Body.String(), expected)
	}
}

func TestPostMessageErrorWhenMessageSent(t *testing.T) {
	requestBodyStr := []byte(`{"title": "title", "message": "message"}`)
	req, err := http.NewRequest("POST", "/messages", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	controller := NewMessageSenderController(&fakeMessageSenderError{})
	controller.PostMessage(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got `%v` want `%v`",
			status, http.StatusInternalServerError)
	}

	expected := "Error when sending message\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got `%v` want `%v`",
			rr.Body.String(), expected)
	}
}

func TestPostMessageErrorBadReq(t *testing.T) {
	requestBodyStr := []byte(`{"badTitle": "title", "message": "message"}`)
	req, err := http.NewRequest("POST", "/messages", bytes.NewBuffer(requestBodyStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	controller := NewMessageSenderController(&fakeMessageSenderSuccess{})
	controller.PostMessage(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got `%v` want `%v`",
			status, http.StatusBadRequest)
	}

	expected := "Invalid json received\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got `%v` want `%v`",
			rr.Body.String(), expected)
	}
}

type fakeMessageSenderSuccess struct{}

func (fms *fakeMessageSenderSuccess) SendMessage(_ string, _ string) error {
	return nil
}

type fakeMessageSenderError struct{}

func (fms *fakeMessageSenderError) SendMessage(_ string, _ string) error {
	return fmt.Errorf("some error\n")
}
