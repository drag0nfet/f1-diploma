package server

import (
	"diploma/internal/handlers"
	"diploma/internal/handlers/bar"
	"diploma/internal/handlers/forum"
	"diploma/internal/handlers/index"
	"diploma/internal/handlers/topic"
	"diploma/internal/handlers/userPage"
	"diploma/internal/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run() {
	router := mux.NewRouter()

	// API маршруты
	{
		// Общие маршруты
		router.HandleFunc("/delete-message/{messageId}", handlers.DeleteMessage)
		router.HandleFunc("/check-auth", handlers.CheckAuth)

		// Страница пользователя
		router.HandleFunc("/logout", userPage.Logout)
		router.HandleFunc("/get-blocks/{username}", userPage.GetBlocks)
		router.HandleFunc("/get-messages-list", userPage.GetMessagesList)

		// Идентификация пользователя
		router.HandleFunc("/register", index.Register)
		router.HandleFunc("/login", index.Login)

		// Работа на странице форума
		router.HandleFunc("/create-topic", forum.CreateTopic)
		router.HandleFunc("/get-topics", forum.GetTopics)
		router.HandleFunc("/delete-topic", forum.DeleteTopic)

		// Работа на странице топика
		router.HandleFunc("/get-topic/{topicId}", topic.GetTopic)
		router.HandleFunc("/get-messages/{topicId}", topic.GetMessages)
		router.HandleFunc("/send-message", topic.SendMessage)
		router.HandleFunc("/block-user/{messageId}", topic.BlockUser)

		// Работа на странице бара
		router.HandleFunc("/get-dishes", bar.GetDishes)
		router.HandleFunc("/delete-dish", bar.DeleteDish)

		// Работа на странице добавления блюда
		router.HandleFunc("/create_dish", bar.CreateDish)
	}

	// Страничные маршруты
	{
		// Детекция авторизации
		router.HandleFunc("/", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/index.html")
		}))
		router.HandleFunc("/web/forum", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/forum.html")
		}))
		router.HandleFunc("/web/bar", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/bar.html")
		}))
		router.HandleFunc("/web/bar/create_dish", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/create_dish.html")
		}))

		// Блокировка неавторизованных
		router.HandleFunc("/account/{username}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/userPage.html")
		}))
		router.HandleFunc("/forum/{topicId}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/topic.html")
		}))
	}

	// Общий маршрут, поэтому в самом низу по порядку регистрации перехода
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Обработчики для статических файлов (JS, CSS)
		if r.URL.Path == "/web/styles/general.css" {
			w.Header().Set("Content-Type", "text/css")
		}
		if r.URL.Path == "/web/js/main.js" || r.URL.Path == "/web/js/forum/createTopic.js" {
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
