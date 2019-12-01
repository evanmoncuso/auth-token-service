package rest

import (
	"encoding/json"
	"net/http"
)

func handleHTTPError(res http.ResponseWriter, err string, statusCode int) {
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	encoder.Encode(map[string]interface{}{
		"errors": err,
	})

	return
}
