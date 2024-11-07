package response

import (
	"net/http"

	"github.com/goccy/go-json"
)

type ApiResponse struct {
	status int
	data   interface{}
}

func (res ApiResponse) Status(status int) ApiResponse {
	res.status = status
	return res
}

func (res ApiResponse) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.status)
	return json.NewEncoder(w).Encode(res.data)
}
