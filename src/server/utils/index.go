package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

const AuthCookieName = "Authenticated-User"

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
