package response

import (
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"net/http"
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
	log.Debug().Int("status", res.status).Msg("response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.status)
	return json.NewEncoder(w).Encode(res.data)
}
