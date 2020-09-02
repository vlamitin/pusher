package messages

import "time"

type Message struct {
	Id         int
	SendStatus int
	SendTime   time.Time
}

func NewMessage(id int, sendStatus int, sendTime time.Time) *Message {
	return &Message{
		Id:         id,
		SendStatus: sendStatus,
		SendTime:   sendTime,
	}
}

type StatusCount struct {
	SendStatus int
	Count      int
}

func NewStatusCount(sendStatus int, count int) *StatusCount {
	return &StatusCount{
		SendStatus: sendStatus,
		Count:      count,
	}
}
