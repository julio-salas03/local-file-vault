package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"local-file-vault/api"
	"local-file-vault/errorcodes"
	"local-file-vault/utils"
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

func GetUserFromAuthCookie(cookieValue string) (string, error) {
	token, err := jwt.Parse(cookieValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return base64.StdEncoding.DecodeString(GetJWTSecret())
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return "", fmt.Errorf("malformed/invalid auth token")
	}

	subject, ok := claims["sub"].(string)

	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid subject in auth token")
	}

	return subject, nil
}

func BuildUserData(user string) map[string]interface{} {
	return map[string]interface{}{
		"username":     user,
		"uploadFolder": user,
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		response := api.Response{
			ErrorCode: errorcodes.BadRequest,
			Message:   "Malformed/Invalid form data",
		}
		w.WriteHeader(http.StatusBadRequest)
		api.WriteResponse(w, response)
		return
	}

	username := strings.TrimSpace(r.Form.Get("username"))

	if username == "" {
		response := api.Response{
			ErrorCode: errorcodes.BadRequest,
			Message:   "Missing 'username' in request data",
		}
		w.WriteHeader(http.StatusBadRequest)
		api.WriteResponse(w, response)
		return
	}

	var dbUser string
	var dbPassword string
	var salt string

	err = utils.Query(func(conn *pgx.Conn) error {
		return conn.QueryRow(context.Background(), "select username, salt, password from users where username=$1", username).Scan(&dbUser, &salt, &dbPassword)
	})

	if err != nil {
		api.InternalServerError(w)
		return
	}

	password := strings.TrimSpace(r.Form.Get("password"))

	if password == "" {
		response := api.Response{
			ErrorCode: errorcodes.BadRequest,
			Message:   "Missing 'password' in request data",
		}
		w.WriteHeader(http.StatusBadRequest)
		api.WriteResponse(w, response)
		return
	}

	argon := argon2.DefaultConfig()

	saltBytes, err := hex.DecodeString(salt)

	if err != nil {
		api.InternalServerError(w)
		return
	}

	encoded, err := argon.Hash([]byte(password), saltBytes)

	if err != nil {
		api.InternalServerError(w)
		return
	}

	hashedPassword := string(encoded.Encode()[:])

	if subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(dbPassword)) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		response := api.Response{
			ErrorCode: errorcodes.InvalidCredentials,
			Message:   "Invalid username or password",
		}
		api.WriteResponse(w, response)
		return
	}

	token, err := GenerateAuthenticationJWT(dbUser)

	if err != nil {
		api.InternalServerError(w)
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

	response := api.Response{
		Message: "Authenticated",
		Data:    BuildUserData(dbUser),
	}

	http.SetCookie(w, &cookie)

	api.WriteResponse(w, response)
}

func HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	authCookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := api.Response{
			ErrorCode: errorcodes.Unauthorized,
			Message:   fmt.Sprintf("Missing '%s' cookie in request header", AuthCookieName),
		}
		api.WriteResponse(w, response)
		return
	}

	username, err := GetUserFromAuthCookie(authCookie.Value)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := api.Response{
			ErrorCode: errorcodes.BadJWT,
			Message:   err.Error(),
		}
		api.WriteResponse(w, response)
		return
	}

	response := api.Response{
		Message: "Authorized",
		Data:    BuildUserData(username),
	}

	w.WriteHeader(http.StatusOK)
	api.WriteResponse(w, response)
}
