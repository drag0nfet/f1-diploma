import { showNotification } from "../notification";

export async function loadBooking() {
    const bookingContainer = document.getElementById("booking-container");
    const bookingBlock = document.querySelector(".booking-block");

    try {
        // Запрос к серверу для получения активных броней пользователя
        const response = await fetch(`/get-bookings`,
            {headers: {'X-Requested-With': 'XMLHttpRequest'}});
        if (!response.ok) {
            throw new Error("Ошибка при загрузке броней");
        }
        let bookings = await response.json();

        if (bookings.length === 0) {
            bookingContainer.innerHTML = "<p>Активные брони отсутствуют.</p>";
            bookingBlock.style.display = "block";
            return;
        }

        bookings = bookings.data

        // Группировка броней по ивенту и залу
        const groupedBookings = bookings.reduce((acc, booking) => {
            const key = `${booking.event.event_id}_${booking.hall.hall_id}`;
            if (!acc[key]) {
                acc[key] = {
                    event: booking.event,
                    hall: booking.hall,
                    spots: []
                };
            }
            acc[key].spots.push({
                table_name: booking.table.table_name,
                spot_name: booking.spot.spot_name,
                booking_id: booking.booking_id
            });
            return acc;
        }, {});

        // Очистка контейнера
        bookingContainer.innerHTML = "";

        // Создание плиток для каждой группы
        Object.values(groupedBookings).forEach(group => {
            const tile = document.createElement("div");
            tile.classList.add("booking-tile");

            // Форматирование даты
            const eventDate = new Date(group.event.time_start).toLocaleString("ru-RU", {
                year: "numeric",
                month: "long",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit"
            });

            // Создание HTML для плитки
            tile.innerHTML = `
                <h3 class="booking-event-title">${group.event.description}</h3>
                <p class="booking-event-date">Дата: ${eventDate}</p>
                <p class="booking-hall-name">Зал: ${group.hall.name}</p>
                <div class="booking-spots-list">
                    ${group.spots.map(spot => `
                        <div class="booking-spot-item">
                            <span>Стол ${spot.table_name} Место ${spot.spot_name}</span>
                            <button class="cancel-booking-btn" data-booking-id="${spot.booking_id}">Отменить бронирование</button>
                        </div>
                    `).join("")}
                </div>
            `;

            bookingContainer.appendChild(tile);
        });

        // Показываем блок бронирований
        bookingBlock.style.display = "block";

        // Обработчик для кнопок отмены
        document.querySelectorAll(".cancel-booking-btn").forEach(button => {
            button.addEventListener("click", async (e) => {
                const bookingId = e.currentTarget.dataset.bookingId;
                if (confirm("Вы уверены, что хотите отменить бронирование?")) {
                    try {
                        const cancelResponse = await fetch(`/cancel-booking?booking_id=${bookingId}`, {
                            method: "POST",
                            headers: {
                                "Content-Type": "application/json",
                                'X-Requested-With': 'XMLHttpRequest'
                            }
                        });
                        if (!cancelResponse.ok) {
                            throw new Error("Ошибка при отмене бронирования");
                        }
                        showNotification("success", "Бронирование успешно отменено");
                        // Перезагружаем брони
                        await loadBooking();
                    } catch (error) {
                        showNotification("error", "Ошибка при отмене бронирования");
                    }
                }
            });
        });
    } catch (error) {
        console.error(error);
        showNotification("error", "Ошибка при загрузке броней");
    }
}