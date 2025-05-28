import {getUnMdText} from "./getUnMdText";

const eventsContainer = document.getElementById("events-container");

export function loadEvents(isFiltered) {
    eventsContainer.innerHTML = '<div class="spinner">Загрузка...</div>';

    const calendarInput = document.querySelector(".calendar-input");
    const categorySelect = document.querySelector(".filter-category select");

    // Собираем данные фильтров
    let dateRange = calendarInput.value.trim();
    const selectedCategories = Array.from(categorySelect.selectedOptions).map(option => option.value);

    // Формируем параметры запроса
    const params = new URLSearchParams();

    // Обработка диапазона дат
    if (!dateRange && isFiltered) {
        // Если диапазон пуст и фильтр применён, используем текущую дату
        const today = new Date();
        const formattedToday = `${String(today.getDate()).padStart(2, "0")}.${String(today.getMonth() + 1).padStart(2, "0")}.${today.getFullYear()}`;
        params.append("date_from", formattedToday);
    } else if (dateRange) {
        if (dateRange.includes("—")) {
            params.append("date_range", dateRange.replace(/\s+/g, ""));
        } else {
            params.append("date_from", dateRange);
        }
    }

    // Обработка категории спорта
    if (selectedCategories.length > 0 && isFiltered) {
        params.append("sport_category", selectedCategories.join(","));
    }

    // Отправка GET-запроса
    fetch(`/get-events-list?${params.toString()}`, {
        method: "GET",
        headers: {
            "X-Requested-With": "XMLHttpRequest"
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success && Array.isArray(data.events)) {
                eventsContainer.innerHTML = '';
                data.events.forEach(event => {
                    const eventElement = document.createElement("div");
                    eventElement.className = "event-item";

                    // Форматируем дату и продолжительность
                    const eventDate = new Date(event.time_start);
                    const formattedDate = eventDate.toLocaleDateString("ru-RU", {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit"
                    });
                    const durationText = event.duration ? `${event.duration} мин` : "Не указано";

                    // Проверяем наличие полей для простой версии ответа
                    const description = getUnMdText(event.description) || "Описание отсутствует";
                    const sportCategory = event.sport_category || "Не указано";
                    const sportType = event.sport_type || "Не указано";

                    eventElement.innerHTML = `
                        <div class="event-content">
                            <div class="event-info">
                                <h3 class="event-title">${description}</h3>
                                <div class="event-details">
                                    <span class="event-category">
                                        <span class="category-dot"></span>
                                        ${sportCategory}
                                    </span>
                                    <span class="event-type"> — ${sportType}</span>
                                </div>
                            </div>
                            <div class="event-meta">
                                <span class="event-date">${formattedDate}</span>
                                <span class="event-duration">${durationText}</span>
                            </div>
                        </div>
                    `;

                    // Добавляем обработчик клика для перехода на страницу ивента
                    eventElement.addEventListener("click", () => {
                        fetch(`/booking/event/${event.event_id}`, {
                            method: "GET",
                            headers: {
                                "X-Requested-With": "XMLHttpRequest"
                            }
                        })
                            .then(response => {
                                if (response.ok) {
                                    window.location.href = `/booking/event/${event.event_id}`;
                                } else {
                                    throw new Error("Ошибка доступа к странице ивента");
                                }
                            })
                            .catch(error => {
                                console.error("Ошибка перехода:", error);
                            });
                    });

                    eventsContainer.appendChild(eventElement);
                });
            } else {
                eventsContainer.innerHTML = '<p>Не удалось загрузить события.</p>';
                console.error("Ошибка: " + data.message);
            }
        })
        .catch(error => {
            console.error("Ошибка загрузки событий:", error);
            eventsContainer.innerHTML = '<p>Ошибка загрузки событий.</p>';
        });
}