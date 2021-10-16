package notifications

import "gohooks/pkg/executor"

type Service interface {
	Send(msg string) error
}

const (
	notifySend = "notify-send"
)

func NewNotifySendService() (Service, error) {
	e, err := executor.New(notifySend)
	if err != nil {
		return nil, err
	}

	return &notifySendService{
		executor: e,
	}, nil
}

// Uses `notify-send` to send notifications
type notifySendService struct {
	executor executor.Executor
}

func (service *notifySendService) Send(msg string) error {
	return service.executor.Run([]string{msg})
}
