import { initSendMessage } from "./sendMessage.js";
import { addMessageToDOM } from "./addMessageToDOM.js";
import { initReplyBtn } from "./replyBtn.js";

async function getUserRights() {
    try {
        const response = await fetch('/check-auth', {
            method: 'GET',
            credentials: 'include',
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        if (response.status === 401) {
            console.warn("Не авторизован");
            return null;
        }

        const data = await response.json();

        if (data.success) {
            return parseInt(data.rights, 10);
        } else {
            console.warn("Не удалось получить rights:", data.message);
            return null;
        }
    } catch (error) {
        console.error("Ошибка при получении прав пользователя:", error);
        return null;
    }
}

export async function loadTopicData(topicId) {
    const rights = await getUserRights();
        fetch(`/get-topic/${topicId}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById("topic-title").textContent = data.topic.title;
                document.querySelector(".header").textContent = data.topic.title;
                document.getElementById("message-form").style.display = "block";
                loadMessages(topicId, rights);
                initSendMessage(topicId, rights);
            } else {
                document.getElementById("topic-title").textContent = "Ошибка: " + data.message;
            }
        })
        .catch(error => console.error("Ошибка загрузки темы:", error));
}

function loadMessages(topicId, rights) {
    fetch(`/get-messages/${topicId}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            const messagesContainer = document.getElementById("messages-container");
            messagesContainer.innerHTML = "";
            if (data.success && Array.isArray(data.messages)) {
                data.messages.forEach(message => {
                    addMessageToDOM(message, rights);
                });
                initReplyBtn();
            }
        })
        .catch(error => console.error("Ошибка загрузки сообщений:", error));
}