package server

import (
	"diploma/internal/handlers"
	"diploma/internal/handlers/account"
	"diploma/internal/handlers/forum"
	"diploma/internal/handlers/index"
	"diploma/internal/handlers/topic"
	"diploma/internal/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run() {
	router := mux.NewRouter()
	//mux := http.NewServeMux()

	// API маршруты
	{
		// Идентификация пользователя
		router.HandleFunc("/register", index.Register)
		router.HandleFunc("/login", index.Login)
		router.HandleFunc("/logout", account.Logout)
		router.HandleFunc("/check-auth", handlers.CheckAuth)

		// Работа на странице форума
		router.HandleFunc("/create-discuss", forum.CreateChat)
		router.HandleFunc("/get-topics", forum.GetTopics)
		router.HandleFunc("/delete-discuss", forum.DeleteDiscuss)

		// Работа на странице топика
		router.HandleFunc("/get-topic/{topicId}", topic.GetTopic)
		router.HandleFunc("/get-messages/{topicId}", topic.GetMessages)
		router.HandleFunc("/send-message", topic.SendMessage)
		router.HandleFunc("/delete-message/{messageId}", handlers.DeleteMessage)
	}

	// Страничные маршруты
	{
		// Детекция авторизации
		router.HandleFunc("/", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/index.html")
		}))
		router.HandleFunc("/web/discuss", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/forum.html")
		}))

		// Блокировка неавторизованных
		router.HandleFunc("/account", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/account.html")
		}))
		router.HandleFunc("/discuss/{topicId}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/topic.html")
		}))
	}

	// Общий маршрут, поэтому в самом низу по порядку регистрации перехода
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Обработчики для статических файлов (JS, CSS)
		if r.URL.Path == "/web/styles/general.css" {
			w.Header().Set("Content-Type", "text/css")
		}
		if r.URL.Path == "/web/js/main.js" || r.URL.Path == "/web/js/forum/createTheme.js" {
			w.Header().Set("Content-Type", "application/javascript")
		}
		fileServer := http.FileServer(http.Dir("web"))
		fileServer.ServeHTTP(w, r)
	})))

	handler := services.EnableCORS(router)

	err := http.ListenAndServe(SetIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
