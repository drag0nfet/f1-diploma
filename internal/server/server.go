package server

import (
	"diploma/internal/handlers"
	"log"
	"net"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web"))
	mux.HandleFunc("/web/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/web/styles.css" {
			w.Header().Set("Content-Type", "text/css")
		}
		if r.URL.Path == "/web/js/script.js" {
			w.Header().Set("Content-Type", "application/javascript")
		}
		http.StripPrefix("/web/", fileServer).ServeHTTP(w, r)
	})

	// API маршруты
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.Register(w, r)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(w, r)
	})
	// Главная страница
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Index(w, r)
	})

	handler := enableCORS(mux)

	err := http.ListenAndServe(setIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	} else {
		log.Println("ListenAndServe success")
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Interfaces error:", err)
		return ""
	}
	for _, iface := range interfaces {
		if iface.Name == "Wi-Fi" || iface.Name == "Беспроводная сеть" {
			addrs, err := iface.Addrs()
			if err != nil {
				log.Println("Wi-Fi search error:", err)
				return ""
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					log.Print("Попытка запуска сервера на адресе http://", ip, ":5051/")
					return ip
				}
			}
		}
	}
	log.Println("Не найден необходимый адрес, сервер не будет запущен.")
	return ""
}
