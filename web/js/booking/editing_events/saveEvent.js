export function saveEvent() {
    const eventId = sessionStorage.getItem("event_id");
    const description = document.getElementById("event-description").value;
    const timeStart = document.getElementById("event-time-start").value;
    const duration = document.getElementById("event-duration").value;
    const sportCategory = document.getElementById("event-sport-category").value;
    const sportType = document.getElementById("event-sport-type").value;
    const priceStatus = document.getElementById("event-price-status").value;

    if (!description || !timeStart || !duration || !sportCategory || !sportType || !priceStatus) {
        alert("Пожалуйста, заполните все обязательные поля.");
        return;
    }

    const dateTimeRegex = /^(\d{2})\.(\d{2})\.(\d{4})\s(\d{2}):(\d{2})$/;
    if (!dateTimeRegex.test(timeStart)) {
        alert("Некорректный формат даты и времени. Используйте формат дд.мм.гггг чч:мм");
        return;
    }

    const [, day, month, year, hours, minutes] = timeStart.match(dateTimeRegex);
    const parsedDate = new Date(year, month - 1, day, hours, minutes);
    if (isNaN(parsedDate.getTime())) {
        alert("Некорректная дата или время.");
        return;
    }

    const eventData = {
        event_id: parseInt(eventId),
        description: description,
        time_start: parsedDate.toISOString(),
        duration: parseInt(duration),
        sport_category: sportCategory,
        sport_type: sportType,
        price_status: priceStatus
    };

    fetch("/save-event", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-Requested-With": "XMLHttpRequest"
        },
        body: JSON.stringify(eventData)
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert("Ивент успешно сохранён!");
                sessionStorage.setItem("event_id", data.event_id);
                window.location.href = "/editing_events";
            } else {
                alert("Ошибка сохранения ивента: " + (data.message || "Неизвестная ошибка"));
            }
        })
        .catch(error => {
            console.error("Ошибка при сохранении ивента:", error);
            alert("Не удалось сохранить ивент: " + error.message);
        });
}