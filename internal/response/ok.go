package response

import "net/http"

func Ok(data interface{}) ApiResponse {
	return ApiResponse{status: http.StatusOK, data: data}
}
