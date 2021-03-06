package messages

import "time"

type Notifier interface {
	Notify(title string, message string) (status int, err error)
}

type MessageRepo interface {
	Save(sendStatus int, sendTime time.Time) error
	GetCountsByTimeFrame(from time.Time, to time.Time) ([]StatusCount, error)
}

type MessageService struct {
	repo     MessageRepo
	notifier Notifier
}

func NewMessageService(repo MessageRepo, notifier Notifier) MessageService {
	return MessageService{repo: repo, notifier: notifier}
}

func (mr *MessageService) SendMessage(title string, message string) error {
	requestTime := time.Now()
	status, err := mr.notifier.Notify(title, message)
	if err != nil {
		return err
	}

	err = mr.repo.Save(status, requestTime)
	if err != nil {
		return err
	}

	return nil
}

func (mr *MessageService) GetCountsByTimeFrame(from time.Time, to time.Time) ([]StatusCount, error) {
	return mr.repo.GetCountsByTimeFrame(from, to)
}
