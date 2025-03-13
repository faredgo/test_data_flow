package req

import (
	"net/http"
	"test_data_flow/pkg/res"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		res.ReturnError(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}
