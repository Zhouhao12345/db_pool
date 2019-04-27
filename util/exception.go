package util

type QueueOutFlowError struct {
	Message string
	Level string
}

func (e QueueOutFlowError) Error() string {
	return e.Message
}
