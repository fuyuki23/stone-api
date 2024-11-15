package response

import "net/http"

func Ok(data interface{}) APIResponse {
	return APIResponse{status: http.StatusOK, data: data}
}
