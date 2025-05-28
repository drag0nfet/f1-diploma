import {loadTopic} from "./loadTopic.js";
import {showNotification} from "../notification";

export function initCreateTheme() {
    const title = document.getElementById("topic_title").value.trim();
    if (!title) {
        showNotification(
            "error",
            `Введите название темы`
        )

        return;
    }

    fetch('/create-topic', {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: JSON.stringify({ title: title })
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {showNotification(
                "success",
                `Тема создана!`
            )

                document.getElementById("topic_title").value = ""; // Очистка поля
                loadTopic(data.topicId, title, true, true); // Добавляем новую тему
            } else {showNotification(
                "error",
                `Ошибка: ${data.message}`
            )

            }
        })
        .catch(error => console.error("Ошибка создания темы:", error));
}
