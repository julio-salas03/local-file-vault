package api

import (
	"encoding/json"
	"local-file-vault/errorcodes"
	"net/http"
)

type Response struct {
	ErrorCode string                 `json:"errorCode,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Message   string                 `json:"message"`
}

func WriteResponse(w http.ResponseWriter, response Response) {
	responseType := "success"

	if response.ErrorCode != "" {
		responseType = "error"
	}

	wrappedResponse := struct {
		Response
		Type string `json:"type"`
	}{
		Response: response,
		Type:     responseType,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wrappedResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func InternalServerError(w http.ResponseWriter) {
	response := Response{
		ErrorCode: errorcodes.InternalServerError,
		Message:   "Internal server error",
	}
	w.WriteHeader(http.StatusInternalServerError)
	WriteResponse(w, response)
}
