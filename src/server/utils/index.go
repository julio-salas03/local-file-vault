package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const AuthCookieName = "JWT-Auth"

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

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")

	if err != nil {
		fmt.Println("Error retrieving file from form-data")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Upload time: %v\n", time.Now())
	fmt.Printf("File Name: %v\n", handler.Filename)
	fmt.Printf("File Size: %v\n", handler.Size)
	fmt.Printf("MIME  : %v\n", handler.Header.Get("Content-Type"))

	tempFile, err := os.Create(filepath.Join("uploads/shared", filepath.Base(handler.Filename)))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	tempFile.Write(fileBytes)
	fmt.Fprint(w, "Successfully uploaded file")
}
