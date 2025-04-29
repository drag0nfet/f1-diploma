package server

import (
	"log"
	"os/exec"
)

func buildFrontend() {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = "web"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Ошибка сборки фронтенда: %v\n%s", err, output)
	}
	log.Println("Фронтенд собран успешно.")
}
