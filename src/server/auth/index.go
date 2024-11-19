package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"local-file-vault/db"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/matthewhartstonge/argon2"
)

const AuthCookieName = "JWT-Auth"

func GetAuthTokenExpireTime() time.Time {
	return time.Now().Add(time.Hour * 24 * 30)
}

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Environment variable required JWT_SECRET is not set")
	}
	return secret
}

func GenerateAuthenticationJWT(user string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user,
		"exp": GetAuthTokenExpireTime(),
	})

	secretKey, err := base64.StdEncoding.DecodeString(GetJWTSecret())

	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	username := strings.TrimSpace(r.Form.Get("username"))

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request. Need to provide a username to authenticate"))
		return
	}

	var dbUser string
	var dbPassword string
	var salt string

	err = db.Query(func(conn *pgx.Conn) error {
		return conn.QueryRow(context.Background(), "select username, salt, password from users where username=$1", username).Scan(&dbUser, &salt, &dbPassword)
	})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	password := strings.TrimSpace(r.Form.Get("password"))

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request. Need to provide a password to authenticate"))
		return
	}

	argon := argon2.DefaultConfig()

	saltBytes, err := hex.DecodeString(salt)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	encoded, err := argon.Hash([]byte(password), saltBytes)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	hashedPassword := string(encoded.Encode()[:])

	if subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(dbPassword)) != 1 {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateAuthenticationJWT(dbUser)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	expirationTime := GetAuthTokenExpireTime()

	cookie := http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expirationTime,
		MaxAge:   int(time.Until(expirationTime).Seconds()),
	}

	http.SetCookie(w, &cookie)
	w.Write([]byte("Logged in"))
}
