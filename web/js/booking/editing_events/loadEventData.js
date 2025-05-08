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
                timeStartInput.value = event.time_start ? new Date(event.time_start).toISOString().slice(0, 16) : ""; // Формат для input[type="datetime-local"]
                durationInput.value = event.duration || "90";
                sportCategorySelect.value = event.sport_category || "";
                sportTypeInput.value = event.sport_type || "";
                sportTypeInput.disabled = !event.sport_category; // Включаем, если категория выбрана
                priceStatusInput.value = event.price_status || "";

            } else {
                console.error("Ошибка загрузки данных ивента:", data.message);
            }
        })
        .catch(error => console.error("Ошибка загрузки данных ивента:", error));
}
