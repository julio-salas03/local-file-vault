package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request)  {
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

func ServeOptimizedFile(filename string, w http.ResponseWriter, r *http.Request) {
    acceptedEncodings := r.Header.Values("Accept-Encoding")
    var filepath strings.Builder
    filepath.WriteString(filename)
    
    for i := 0; i < len(acceptedEncodings); i++ {
        if strings.Contains(acceptedEncodings[i],"br") {
            w.Header().Add("Content-Encoding","br")
            filepath.WriteString(".br")
            break
        } 
    }

    bytes, err := os.ReadFile(filepath.String())
    
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w,"Internal Server Error")
        return
    }
    
    w.Write(bytes)
}

func main() {
    http.HandleFunc("/api/upload", UploadHandler)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type","text/html; charset=utf-8")
        ServeOptimizedFile("./dist/index.html",w,r)
    })

    http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type","text/javascript;charset=UTF-8")
        ServeOptimizedFile("./dist/index.js",w,r)
    })

    http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type","text/css; charset=utf-8")
        ServeOptimizedFile("./dist/index.css",w,r)
    })

    http.HandleFunc("/Inter-Variable.ttf", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type","font/ttf")
        ServeOptimizedFile("./dist/Inter-Variable.ttf",w,r)
    })

	portNumber := ":8080"

    log.Println("Started on port", portNumber)

    err := http.ListenAndServe(portNumber, nil)

    if err != nil {
        log.Fatal(err)
    }
}


