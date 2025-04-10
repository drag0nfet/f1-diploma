import {loadTopic} from "./loadTopic.js";

export function loadForumData(isModerator) {
    fetch('/get-topics', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            const topicsContainer = document.getElementById("topics-container");
            topicsContainer.innerHTML = ""; // Очищаем контейнер
            if (data.success && Array.isArray(data.topics)) {
                data.topics.forEach(topic => {
                    loadTopic(topic.chat_id, topic.title, isModerator);
                });
            }
        })
        .catch(error => console.error("Ошибка загрузки тем:", error));
}