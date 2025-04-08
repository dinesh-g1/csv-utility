package types

import (
	"encoding/json"
	"fmt"
)

type Error interface {
	Error() string
	ErrorMessage() ([]byte, error)
	ErrorStatusCode() (int, map[string]string)
}
type ApiError struct {
	Cause      error  `json:"-"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (a *ApiError) Error() string {
	if a.Cause == nil {
		return a.Message
	}
	return a.Cause.Error()
}

func (a *ApiError) ErrorMessage() ([]byte, error) {
	msg, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal error message: %v", err)
	}
	return msg, nil
}

func (a *ApiError) ErrorStatusCode() (int, map[string]string) {
	return a.StatusCode, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}
