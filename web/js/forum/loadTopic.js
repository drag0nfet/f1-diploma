import {deleteTopic} from "./deleteTopic.js";

export function loadTopic(topicId, title, isModerator, prepend = false) {
    const topicsContainer = document.getElementById("topics-container");
    const topicElement = document.createElement("div");
    topicElement.className = "topic-item";

    // Создаём HTML для панельки
    let html = `<a href="/forum/${topicId}" class="topic-link">${title}</a>`;
    if (isModerator) {
        html += `<button class="delete-topic-btn" data-chat-id="${topicId}">✖</button>`;
    }
    topicElement.innerHTML = html;

    // Добавляем элемент в контейнер
    if (prepend) {
        topicsContainer.prepend(topicElement);
    } else {
        topicsContainer.appendChild(topicElement);
    }

    // Привязываем обработчик к кнопке удаления (если она есть)
    if (isModerator) {
        const deleteBtn = topicElement.querySelector(".delete-topic-btn");
        deleteBtn.addEventListener("click", deleteTopic);
    }
}