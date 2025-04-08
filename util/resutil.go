package util

import (
	"fmt"
	"net/http"
)

func SendResponse(w http.ResponseWriter, response *string) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, *response)
}
