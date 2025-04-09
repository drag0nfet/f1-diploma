package services

import (
	"diploma/internal/database"
	"diploma/internal/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// CheckAuthCookie returns username, id, rights, response
func CheckAuthCookie(r *http.Request) (string, int, int, Response) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		return "", 0, 0, Response{Success: false, Message: "Не авторизован"}
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("Wv1%`j9pr]0d[s'_HwX,U|m;6^3>u="), nil
	})
	if err != nil || !token.Valid {
		return "", 0, 0, Response{Success: false, Message: "Неверный токен"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", 0, 0, Response{Success: false, Message: "Ошибка чтения токена"}
	}

	userIDFloat, _ := claims["user_id"].(float64) // JWT хранит числа как float64
	userID := int(userIDFloat)

	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return "", 0, 0, Response{Success: false, Message: "Пользователь не найден"}
	}

	return user.Login, user.UserID, user.Rights, Response{Success: true}
}
