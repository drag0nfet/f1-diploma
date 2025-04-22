package server

import (
	"diploma/internal/handlers"
	"diploma/internal/handlers/bar"
	"diploma/internal/handlers/forum"
	"diploma/internal/handlers/index"
	"diploma/internal/handlers/index/news"
	"diploma/internal/handlers/topic"
	"diploma/internal/handlers/userPage"
	"diploma/internal/handlers/userPage/moderatorBlocks"
	"diploma/internal/handlers/userPage/userBlocks"
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

		// Страница домашняя
		router.HandleFunc("/register", index.Register)
		router.HandleFunc("/login", index.Login)

		// Редактирование новости
		router.HandleFunc("/update-news", news.UpdateNews)
		router.HandleFunc("/delete-news/{news_id}", news.DeleteNews)
		router.HandleFunc("/create-news", news.CreateNews)
		router.HandleFunc("/load-news-by-status", news.LoadNews)
		router.HandleFunc("/load-news-info", news.LoadNewsInfo)

		// Страница пользователя
		router.HandleFunc("/logout", userPage.Logout)
		router.HandleFunc("/get-messages-list", userPage.GetMessagesList)
		router.HandleFunc("/get-blocks/{username}", userBlocks.GetBlocks)
		router.HandleFunc("/submit-unblock-request", userBlocks.UnblockRequest)
		router.HandleFunc("/get-requests", moderatorBlocks.GetRequests)
		router.HandleFunc("/approve/{request_id}", moderatorBlocks.Approve)
		router.HandleFunc("/reject/{request_id}", moderatorBlocks.Reject)

		// Страница форума
		router.HandleFunc("/create-topic", forum.CreateTopic)
		router.HandleFunc("/get-topics", forum.GetTopics)
		router.HandleFunc("/delete-topic", forum.DeleteTopic)

		// Страница топика
		router.HandleFunc("/get-topic/{topicId}", topic.GetTopic)
		router.HandleFunc("/get-messages/{topicId}", topic.GetMessages)
		router.HandleFunc("/send-message", topic.SendMessage)
		router.HandleFunc("/block-user/{messageId}", topic.BlockUser)

		// Страница бара
		router.HandleFunc("/get-dishes", bar.GetDishes)
		router.HandleFunc("/delete-dish", bar.DeleteDish)

		// Страница добавления блюда
		router.HandleFunc("/create-dish", bar.CreateDish)
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
		router.HandleFunc("/account/{username}", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/userPage.html")
		}))
		router.HandleFunc("/news-list", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/news-list.html")
		}))
		router.HandleFunc("/editing_news", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/editing_news.html")
		}))

		// Блокировка неавторизованных
		router.HandleFunc("/forum/{topicId}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/pages/topic.html")
		}))
	}

	// Общий маршрут, поэтому в самом низу по порядку регистрации перехода
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	handler := services.EnableCORS(router)

	err := http.ListenAndServe(SetIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
