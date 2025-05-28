import {loadTables} from "./loadTables";

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
        } else {
            alert(data.message || "Ошибка бронирования");
            console.error("Ошибка:", data.message);
        }
    } catch (error) {
        console.error("Ошибка при бронировании/отмене:", error);
        alert("Произошла ошибка при выполнении действия");
    }
}

export function getEventIdFromUrl() {
    const pathParts = window.location.pathname.split("/");
    return pathParts[pathParts.length - 1];
}