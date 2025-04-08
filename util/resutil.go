package util

import (
	"encoding/json"
	"github.com/dinesh-g1/csv-utility/types"
	"net/http"
)

func SendResponse(w http.ResponseWriter, response *types.SuccessResponse) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return err
	}
	return nil
}
