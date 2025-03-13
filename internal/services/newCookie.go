package services

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func NewCookie(w http.ResponseWriter, username string) http.Cookie {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte("Wv1%`j9pr]0d[s'_HwX,U|m;6^3>u="))
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
		MaxAge:   300,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	return cookie
}
