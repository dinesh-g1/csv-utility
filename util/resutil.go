package util

import (
	"fmt"
	"net/http"
)

func SendResponse(w http.ResponseWriter, response *string) error {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, *response)
	if err != nil {
		return err
	}
	return nil
}
