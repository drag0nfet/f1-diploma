import {addTopicToDOM} from "./addTopicToDOM.js";

export function initCreateTheme(isModerator) {
    const title = document.getElementById("discuss_title").value.trim();
    if (!title) {
        alert("Введите название темы!");
        return;
    }

    fetch('/create-discuss', {
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
            if (data.success) {
                alert("Тема создана!");
                document.getElementById("discuss_title").value = ""; // Очистка поля
                addTopicToDOM(data.topicId, title, isModerator, true); // Добавляем новую тему
            } else {
                alert("Ошибка: " + data.message);
            }
        })
        .catch(error => console.error("Ошибка создания темы:", error));
}
