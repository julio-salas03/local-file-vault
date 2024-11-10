package main

import (
	"context"
	"fmt"
	"io"
	"local-file-vault/db"
	"local-file-vault/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	utils.LoadEnvFile()

	http.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		// TODO: refactor this to something more secure, like a JWT
		err := r.ParseMultipartForm(10 << 20)
		var retrievedUser string

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if len(password) <= 0 || len(username) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "You must provide both a username and password to login")
			return
		}

		err = db.Query(func(conn *pgx.Conn) error {
			return conn.QueryRow(context.Background(), "select username from users where username=$1 and password=$2", username, password).Scan(&retrievedUser)
		})

		if len(retrievedUser) <= 0 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Couldn't find user with the provided username and password")
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error")
			return
		}

		cookie := http.Cookie{
			Name:     utils.AuthCookieName,
			Value:    retrievedUser,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)
		fmt.Fprintf(w, "Authenticated")
	})

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

	var port strings.Builder
	port.WriteString(":")
	port.WriteString(os.Getenv("PORT"))

	log.Println("Started on port", port.String())

	err := http.ListenAndServe(port.String(), nil)

	if err != nil {
		log.Fatal(err)
	}
}
