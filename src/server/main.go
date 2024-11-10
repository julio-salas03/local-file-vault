package main

import (
	"fmt"
	"local-file-vault/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	utils.LoadEnvFile()

	http.HandleFunc("/api/upload", utils.HandleFileUpload)

	http.HandleFunc("/api/login", utils.HandleLogin)

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

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	if port == ":" {
		port = ":8080"
	}

	log.Println("Started on port", port)

	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal(err)
	}
}
