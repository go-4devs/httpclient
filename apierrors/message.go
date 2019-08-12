package apierrors

// Message api error with message
type Message struct {
	Message string `json:"message"`
}

func (m *Message) Error() string {
	return m.Message
}

// ErrorMessage create error
func ErrorMessage() error {
	return &Message{}
}
