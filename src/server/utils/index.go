package utils

import (
	"context"
	"fmt"
	"io"
	"local-file-vault/db"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
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

func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
		Name:     AuthCookieName,
		Value:    retrievedUser,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "Authenticated")
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
