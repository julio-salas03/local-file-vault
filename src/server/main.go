package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func uploadHandler(w http.ResponseWriter, r *http.Request)  {
    r.ParseMultipartForm(10<<20)
    file, handler, err:= r.FormFile("file")

    if err !=  nil {
        fmt.Println("Error retrieving file from form-data")
        fmt.Println(err)
        return
    }

    defer file.Close()

    fmt.Printf("Upload time: %v\n",time.Now())
    fmt.Printf("File Name: %v\n",handler.Filename)
    fmt.Printf("File Size: %v\n",handler.Size)
    fmt.Printf("MIME Header: %v\n",handler.Header)

    tempFile,err := os.Create(filepath.Join("uploads", filepath.Base(handler.Filename)))

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
    fmt.Fprint(w,"Successfully uploaded file")
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("./dist")))
    http.HandleFunc("/api/upload", uploadHandler)

	portNum := ":8080"

    log.Println("Started on port", portNum)

    err := http.ListenAndServe(portNum, nil)

    if err != nil {
        log.Fatal(err)
    }
}


