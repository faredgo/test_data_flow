package res

import (
	"encoding/json"
	"net/http"
)

func ReturnError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: errorMessage,
	})
}
