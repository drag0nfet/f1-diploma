package server

import (
	"diploma/internal/handlers"
	"diploma/internal/handlers/bar"
	"diploma/internal/handlers/booking/editing_news"
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
	buildFrontend()

	router := mux.NewRouter()

	// API маршруты
	{
		// Общие маршруты
		router.HandleFunc("/delete-message/{messageId}", handlers.DeleteMessage)
		router.HandleFunc("/check-auth", handlers.CheckAuth)

		// Страница домашняя
		router.HandleFunc("/register", index.Register)
		router.HandleFunc("/login", index.Login)

		// Создание и редактирование новости
		router.HandleFunc("/update-news", news.UpdateNews)
		router.HandleFunc("/delete-news/{news_id}", news.DeleteNews)
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

		// Страница редактирования мест и залов - модераторская
		router.HandleFunc("/get-halls-list", editing_news.GetHallsList)
		router.HandleFunc("/get-hall", editing_news.GetHall)
		router.HandleFunc("/save-hall", editing_news.SaveHall)
		router.HandleFunc("/delete-table", editing_news.DeleteTable)
		router.HandleFunc("/get-spot-count", editing_news.GetSpotCount)
		router.HandleFunc("/save-table", editing_news.SaveTable)
	}

	// Страничные маршруты
	{
		// Детекция авторизации
		router.HandleFunc("/", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/index.html")
		}))
		router.HandleFunc("/forum", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/forum.html")
		}))
		router.HandleFunc("/bar", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/bar.html")
		}))
		router.HandleFunc("/bar/create_dish", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/create_dish.html")
		}))
		router.HandleFunc("/account/{username}", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/userPage.html")
		}))
		router.HandleFunc("/news-list", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/news-list.html")
		}))
		router.HandleFunc("/editing_news", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/editing_news.html")
		}))
		router.HandleFunc("/booking", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/booking.html")
		}))
		router.HandleFunc("/editing_booking", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/editing_booking.html")
		}))
		router.HandleFunc("/modal_new_table", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/modal_new_table.html")
		}))

		// Блокировка неавторизованных
		router.HandleFunc("/forum/{topicId}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/topic.html")
		}))
	}

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("dir/assets"))))
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("dir"))))

	handler := services.EnableCORS(router)

	err := http.ListenAndServe(SetIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
