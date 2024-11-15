package response

import (
	"net/http"

	"github.com/goccy/go-json"
)

type APIResponse struct {
	status int
	data   interface{}
}

func (res APIResponse) Status(status int) APIResponse {
	res.status = status
	return res
}

func (res APIResponse) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.status)
	return json.NewEncoder(w).Encode(res.data) // nolint:wrapcheck
}
