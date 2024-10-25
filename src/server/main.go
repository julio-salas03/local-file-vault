package main

import (
	"log"
	"net/http"
)

func main() {
    http.Handle("/", http.FileServer(http.Dir("./dist")))

	portNum := ":8080"

    log.Println("Started on port", portNum)

    err := http.ListenAndServe(portNum, nil)
    if err != nil {
        log.Fatal(err)
    }
}


