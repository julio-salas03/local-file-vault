package main

import (
	"fmt"
	"io"
	"local-file-vault/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("MIME Header: %v\n", handler.Header)

	tempFile, err := os.Create(filepath.Join("uploads", filepath.Base(handler.Filename)))

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

func main() {
	utils.LoadEnvFile()

	fmt.Println(os.Getenv("DATABASE_URL"))

	http.HandleFunc("/api/upload", UploadHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		utils.ServeOptimizedFile("./dist/index.html", w, r)
	})

	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript;charset=UTF-8")
		utils.ServeOptimizedFile("./dist/index.js", w, r)
	})

	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css; charset=utf-8")
		utils.ServeOptimizedFile("./dist/index.css", w, r)
	})

	http.HandleFunc("/Inter-Variable.ttf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "font/ttf")
		utils.ServeOptimizedFile("./dist/Inter-Variable.ttf", w, r)
	})

	portNumber := ":8080"

	log.Println("Started on port", portNumber)

	err := http.ListenAndServe(portNumber, nil)

	if err != nil {
		log.Fatal(err)
	}
}
