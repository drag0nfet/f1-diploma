import {updateMarkdownPreview} from "./main";

export function loadEventData(eventId) {
    const descriptionInput = document.getElementById("event-description");
    const timeStartInput = document.getElementById("event-time-start");
    const durationInput = document.getElementById("event-duration");
    const sportCategorySelect = document.getElementById("event-sport-category");
    const sportTypeInput = document.getElementById("event-sport-type");
    const priceStatusInput = document.getElementById("event-price-status");

    const selectedEvent = document.querySelector(".selected-event");
    selectedEvent.style.display = "block";

    sessionStorage.setItem("event_id", eventId);

    if (eventId === -1) {
        descriptionInput.value = "";
        timeStartInput.value = "";
        durationInput.value = "90";
        sportCategorySelect.value = "";
        sportTypeInput.value = "";
        sportTypeInput.disabled = true;
        priceStatusInput.value = "";
        updateMarkdownPreview();
        return;
    }

    fetch(`/get-event?event_id=${encodeURIComponent(eventId)}`, {
        method: 'GET',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success && data.event) {
                const event = data.event;

                descriptionInput.value = event.description || "";
                // Преобразуем time_start из ISO в формат дд.мм.гггг чч:мм
                if (event.time_start) {
                    const date = new Date(event.time_start);
                    const day = String(date.getDate()).padStart(2, "0");
                    const month = String(date.getMonth() + 1).padStart(2, "0");
                    const year = date.getFullYear();
                    const hours = String(date.getHours()).padStart(2, "0");
                    const minutes = String(date.getMinutes()).padStart(2, "0");
                    timeStartInput.value = `${day}.${month}.${year} ${hours}:${minutes}`;
                } else {
                    timeStartInput.value = "";
                }
                durationInput.value = event.duration || "90";
                sportCategorySelect.value = event.sport_category || "";
                sportTypeInput.value = event.sport_type || "";
                sportTypeInput.disabled = !event.sport_category;
                priceStatusInput.value = event.price_status || "";
                updateMarkdownPreview();
            } else {
                console.error("Ошибка загрузки данных ивента:", data.message);
            }
        })
        .catch(error => console.error("Ошибка загрузки данных ивента:", error));
}