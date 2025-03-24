export function loadTopicData(topicId) {
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
                document.getElementById("message-form").style.display = "block";
                loadMessages(topicId);
            } else {
                document.getElementById("topic-title").textContent = "Ошибка: " + data.message;
            }
        })
        .catch(error => console.error("Ошибка загрузки темы:", error));
}

function loadMessages(topicId) {
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
                    const messageElement = document.createElement("div");
                    messageElement.className = "message-item";
                    messageElement.textContent = message.text;
                    messagesContainer.appendChild(messageElement);
                });
            }
        })
        .catch(error => console.error("Ошибка загрузки сообщений:", error));
}