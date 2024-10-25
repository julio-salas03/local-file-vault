package main

import (
	"fmt"
	"log"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("./dist")))
    http.HandleFunc("/api/ping", Ping)
    
	portNum := ":8080"

    log.Println("Started on port", portNum)

    err := http.ListenAndServe(portNum, nil)
    if err != nil {
        log.Fatal(err)
    }
}


