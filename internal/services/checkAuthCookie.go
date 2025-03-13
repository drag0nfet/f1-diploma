package services

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// CheckAuthCookie проверяет наличие и валидность куки "auth" с JWT-токеном
func CheckAuthCookie(r *http.Request) (string, Response) {
	// Извлекаем куки
	cookie, err := r.Cookie("auth")
	if err != nil {
		return "", Response{
			Success: false,
			Message: "Не авторизован: куки отсутствует",
		}
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
		}
	}

	// Проверяем, валиден ли токен
	if !token.Valid {
		return "", Response{
			Success: false,
			Message: "Токен недействителен",
		}
	}

	// Извлекаем claims (данные из токена)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", Response{
			Success: false,
			Message: "Не удалось извлечь данные из токена",
		}
	}

	// Извлекаем username
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return "", Response{
			Success: false,
			Message: "Логин пользователя не найден в токене",
		}
	}

	// Проверяем срок действия
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return "", Response{
			Success: false,
			Message: "Срок действия токена истёк",
		}
	}

	// Возвращаем логин пользователя и успешный ответ
	return username, Response{
		Success: true,
	}
}
