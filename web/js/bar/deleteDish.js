import {loadBarData} from "./loadBarData.js";
import {showNotification} from "../notification";

export function deleteDish(event) {
    const dishId = event.target.getAttribute("data-dish-id");
    if (!dishId) {
        console.error("dish_id не найден");
        return;
    }

    if (!confirm("Вы уверены, что хотите удалить это блюдо?")) {
        return;
    }

    fetch('/delete-dish', {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: JSON.stringify({ dish_id: parseInt(dishId) })
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {showNotification(
                "success",
                "Блюдо удалено!"
            )
                loadBarData(true);
            } else {
                showNotification(
                    "error",
                    `Ошибка ${data.message}`
                )
            }
        })
        .catch(error => console.error("Ошибка удаления блюда:", error));
}