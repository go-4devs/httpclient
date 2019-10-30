package apierrors

import (
	"fmt"
	"net/http"
)

// Message api error with message
type Message struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (m *Message) Error() string {
	return fmt.Sprintf("StatusCode:%d, Message: %s", m.StatusCode, m.Message)
}

// ErrorMessage create error
// deprecated: use HTTPErrorMessage
func ErrorMessage() error {
	return &Message{}
}

// HTTPErrorMessage create error with status code
func HTTPErrorMessage(r *http.Response) error {
	return &Message{
		StatusCode: r.StatusCode,
	}
}
