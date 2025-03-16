package services

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// CheckAuthCookie проверяет наличие и валидность куки "auth" с JWT-токеном
func CheckAuthCookie(r *http.Request) (string, Response, int) {
	// Извлекаем куки
	cookie, err := r.Cookie("auth")
	if err != nil {
		return "", Response{
			Success: false,
			Message: "Не авторизован: куки отсутствует",
		}, 0
	}

	// Парсим JWT-токен
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// Возвращаем секретный ключ
		return []byte("Wv1%`j9pr]0d[s'_HwX,U|m;6^3>u="), nil
	})

	if err != nil {
		return "", Response{
			Success: false,
			Message: "Недействительный токен: " + err.Error(),
		}, 0
	}

	// Проверяем, валиден ли токен
	if !token.Valid {
		return "", Response{
			Success: false,
			Message: "Токен недействителен",
		}, 0
	}

	// Извлекаем claims (данные из токена)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", Response{
			Success: false,
			Message: "Не удалось извлечь данные из токена",
		}, 0
	}

	// Извлекаем username
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return "", Response{
			Success: false,
			Message: "Логин пользователя не найден в токене",
		}, 0
	}

	// Проверяем срок действия
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return "", Response{
			Success: false,
			Message: "Срок действия токена истёк",
		}, 0
	}

	rights, ok := claims["rights"].(float64)
	if !ok {
		return "", Response{
			Success: false,
			Message: "Права пользователя нарушены",
		}, 0
	}

	// Возвращаем логин пользователя, его права и успешный ответ
	return username, Response{
		Success: true,
	}, int(rights)
}
