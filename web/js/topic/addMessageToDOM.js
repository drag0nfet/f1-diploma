export function addMessageToDOM(message, rights) {
    const messagesContainer = document.getElementById("messages-container");
    const messageElement = document.createElement("div");
    messageElement.className = "message-item";

    // Проверяем, является ли пользователь модератором
    const isModerator = (rights % 2 === 1);

    // Проверяем, является ли сообщение ответом
    const isReply = message.reply_id !== undefined && message.reply_id !== null;

    // Добавляем классы для стилизации
    if (isReply) {
        messageElement.classList.add("reply");
    }
    if (isModerator) {
        messageElement.classList.add("moderator");
    }

    // Форматируем время
    const date = new Date(message.message_time);
    const formattedTime = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;

    // Создаём HTML для сообщения
    messageElement.innerHTML = `
        <div class="message-header">
            <div>
                <span class="message-author">${message.username}</span>
                <span class="message-id">#${message.message_id}</span>
            </div>
        </div>
        ${isReply ? `<div class="reply-to">в ответ на <a href="#" 
            class="reply-link" data-message-id="${message.reply_id}">#${message.reply_id}</a></div>` : ''}
        <div class="message-content">${message.value}</div>
        <div class="message-timestamp">${formattedTime}</div>
        <div class="message-actions">
            <button class="reply-btn">Ответить</button>
            <button class="delete-btn">Удалить</button>
            <button class="ban-btn">Заблокировать</button>
        </div>
    `;

    messagesContainer.appendChild(messageElement);

    // Привязываем обработчик для ссылок reply-link
    if (isReply) {
        const replyLink = messageElement.querySelector(".reply-link");
        replyLink.addEventListener("click", function (event) {
            event.preventDefault();
            const targetMessageId = this.getAttribute("data-message-id");
            const targetMessage = Array.from(document.querySelectorAll(".message-id"))
                .find(el => el.textContent === `#${targetMessageId}`)
                ?.closest(".message-item");
            if (targetMessage) {
                targetMessage.scrollIntoView({ behavior: "smooth", block: "center" });
            }
        });
    }

    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}