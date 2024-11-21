package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

func ServeOptimizedFile(filename string, w http.ResponseWriter, r *http.Request) {
	acceptedEncodings := r.Header.Values("Accept-Encoding")
	var filepath strings.Builder
	filepath.WriteString(filename)

	for i := 0; i < len(acceptedEncodings); i++ {
		if strings.Contains(acceptedEncodings[i], "br") {
			w.Header().Add("Content-Encoding", "br")
			filepath.WriteString(".br")
			break
		}
	}

	bytes, err := os.ReadFile(filepath.String())

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}

	w.Write(bytes)
}

func LoadEnvFile() {
	directory, err := os.Getwd()

	if err != nil {
		log.Fatal("Could read working directory. Unable to load .env file")
	}

	godotenv.Load(path.Join(directory, ".env"))
}

type APIResponse struct {
	ErrorCode string                 `json:"errorCode,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Message   string                 `json:"message"`
}

func WriteAPIResponse(w http.ResponseWriter, response APIResponse) {
	responseType := "success"

	if response.ErrorCode != "" {
		responseType = "error"
	}

	wrappedResponse := struct {
		APIResponse
		Type string `json:"type"`
	}{
		APIResponse: response,
		Type:        responseType,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wrappedResponse); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
