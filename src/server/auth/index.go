package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
		"exp": GetAuthTokenExpireTime().Unix(),
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

	response := map[string]string{
		"username":     dbUser,
		"uploadFolder": dbUser,
	}

	http.SetCookie(w, &cookie)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	token, err := jwt.Parse(authCookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return base64.StdEncoding.DecodeString(GetJWTSecret())
	})

	if err != nil {
		fmt.Println(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subject, _ := claims["sub"].(string)

		response := map[string]string{
			"username":     subject,
			"uploadFolder": subject,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("Invalid token claims or token invalid:", err)
		return
	}
}
