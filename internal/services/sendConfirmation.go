package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
)

func SendConfirmation(email string) (string, error) {
	confToken, err := generateToken(20)
	if err != nil {
		return "", errors.New("ошибка при генерации токена подтверждения")
	}

	if err = sendConfirmationEmail(email, confToken); err != nil {
		return "", errors.New("ошибка при отправке письма подтверждения")
	}
	return confToken, nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, (length*3)/4+1)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	if len(token) > length {
		token = token[:length]
	}
	return token, nil
}

func sendConfirmationEmail(toEmail, token string) error {
	ip := os.Getenv("IP")
	gmail := os.Getenv("GMAIL")
	gpass := os.Getenv("GPASS")

	e := email.NewEmail()
	e.From = fmt.Sprintf("F1-Diploma <%s>", gmail)
	e.To = []string{toEmail}
	e.Subject = "Подтверждение регистрации"
	e.Text = []byte(fmt.Sprintf("Перейдите по ссылке для подтверждения вашей электронной почты: %s/confirm?token=%s", ip, token))
	return e.Send("smtp.gmail.com:587", smtp.PlainAuth("", gmail, gpass, "smtp.gmail.com"))
}
