package res

import (
	"encoding/json"
	"net/http"
)

const (
	contentType     string = "Content-Type"
	applicationJson string = "application/json"
)

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
