package grpc

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error    uint32 `json:"error"`
	ErrorMsg string `json:"error_msg"`
}

func sendOkResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, errorType int, errorCode uint32, errorMsg string) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(errorType)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errorCode, ErrorMsg: errorMsg})
}
