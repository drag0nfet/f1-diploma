package server

import (
	"diploma/internal/handlers"
	"diploma/internal/handlers/bar"
	"diploma/internal/handlers/booking"
	"diploma/internal/handlers/booking/editing_events"
	"diploma/internal/handlers/booking/editing_halls"
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
	"time"
)

func Run() {
	start := time.Now()
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
		router.HandleFunc("/get-bookings", userPage.GetBookings)
		router.HandleFunc("/cancel-booking", userPage.CancelBooking)
		router.HandleFunc("/get-booking-pass", userPage.GetBookingPass)

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

		// Страница бронирования спотов
		// router.HandleFunc("/get-events-list", editing_events.GetEventsList) - тот же функционал
		router.HandleFunc("/booking/hall/{hallId}/tables", booking.GetTablesList)
		router.HandleFunc("/booking/hall/{hallId}", booking.GetHallDetails)
		router.HandleFunc("/book-spot", booking.BookSpot)

		// Страница редактирования мест и залов - модераторская
		router.HandleFunc("/get-halls-list", editing_halls.GetHallsList)
		router.HandleFunc("/get-hall", editing_halls.GetHall)
		router.HandleFunc("/save-hall", editing_halls.SaveHall)
		router.HandleFunc("/delete-hall", editing_halls.DeleteHall)
		router.HandleFunc("/get-spot-count", editing_halls.GetSpotCount)
		router.HandleFunc("/save-table", editing_halls.SaveTable)
		router.HandleFunc("/delete-table", editing_halls.DeleteTable)

		// Страница редактирования ивентов - модераторская
		router.HandleFunc("/get-events-list", editing_events.GetEventsList)
		router.HandleFunc("/get-event", editing_events.GetEvent)
		router.HandleFunc("/save-event", editing_events.SaveEvent)
		router.HandleFunc("/delete-event", editing_events.DeleteEvent)
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
		router.HandleFunc("/editing_halls", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/editing_halls.html")
		}))
		router.HandleFunc("/editing_events", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/editing_events.html")
		}))
		router.HandleFunc("/modal_new_table", AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/modal_new_table.html")
		}))

		// Блокировка неавторизованных
		router.HandleFunc("/forum/{topicId}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/topic.html")
		}))
		router.HandleFunc("/booking/event/{event_id}", StrictAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "dir/pages/booking_event.html")
		}))
	}

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("dir/assets"))))
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("dir"))))

	handler := services.EnableCORS(router)

	log.Println("Сервер запущен за ", time.Since(start))
	err := http.ListenAndServe(SetIP()+":5051", handler)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
