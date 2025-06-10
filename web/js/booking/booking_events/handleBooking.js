import { loadTables } from "./loadTables";
import {showNotification} from "../../notification";

export async function handleBooking(tableId, spotId, isBookedByMe) {
    const hallSelect = document.getElementById("hall-select");

    const action = isBookedByMe ? "cancel" : "book";
    try {
        const response = await fetch(`/book-spot`, {
            method: "POST",
            headers: {
                "X-Requested-With": "XMLHttpRequest",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                event_id: getEventIdFromUrl(),
                table_id: tableId,
                spot_id: spotId || null,
                action: action
            })
        });
        const data = await response.json();
        if (data.success) {
            await loadTables(hallSelect.value);
            // Показываем уведомление в зависимости от действия
            showNotification(
                "success",
                action === "book"
                    ? "Вы успешно забронировали место, можете найти пропуск на странице в личном кабинете"
                    : "Бронь места отменена"
            );
        } else {
            showNotification(data.message);
            console.error("Ошибка:", data.message);
        }
    } catch (error) {
        console.error("Ошибка при бронировании/отмене:", error);
        showNotification("Произошла ошибка при выполнении действия");
    }
}

export function getEventIdFromUrl() {
    const pathParts = window.location.pathname.split("/");
    return parseInt(pathParts[pathParts.length - 1], 10);
}