import {addTopicToDOM} from "./createTheme.js";

export function initAuthStatus() {
    const guestContent = document.getElementById("guest-content");
    const userContent = document.getElementById("user-content");
    const moderatorContent = document.getElementById("moderator-content");

    fetch('/check-auth', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success && data.username) {
                const rights = data.rights || 0;

                // Скрываем гостевой контент для всех авторизованных
                guestContent.style.display = "none";

                // Показываем userContent для всех авторизованных
                userContent.style.display = "block";

                // Показываем moderatorContent только для модераторов (rights % 2 === 1)
                if (rights % 2 === 1) {
                    moderatorContent.style.display = "block";
                } else {
                    moderatorContent.style.display = "none";
                }
                loadTopics();
            } else {
                // Не авторизован
                guestContent.style.display = "block";
                userContent.style.display = "none";
                moderatorContent.style.display = "none";
                console.log("Пользователь не авторизован:", data.message);
            }
        })
        .catch(error => {
            console.error("Ошибка проверки авторизации:", error);
            guestContent.style.display = "block";
            userContent.style.display = "none";
            moderatorContent.style.display = "none";
        });
}

function loadTopics() {
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
                    addTopicToDOM(topic.chat_id, topic.title);
                });
            }
        })
        .catch(error => console.error("Ошибка загрузки тем:", error));
}