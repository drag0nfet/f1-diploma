package services

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func NewCookie(w http.ResponseWriter, username string, rights, userId int) http.Cookie {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id":  userId,
		"rights":   rights,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	key := os.Getenv("TOKEN_KEY")

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		response := Response{Success: false, Message: "Ошибка генерации токена"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusInternalServerError)
		var zeroCookie http.Cookie
		return zeroCookie
	}

	cookie := http.Cookie{
		Name:     "auth",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie
}
