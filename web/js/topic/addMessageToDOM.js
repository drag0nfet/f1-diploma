export function addMessageToDOM(message) {
    const messagesContainer = document.getElementById("messages-container");
    const messageElement = document.createElement("div");
    messageElement.className = "message-item";

    // Проверяем, является ли пользователь модератором (заглушка, пока без серверной логики)
    const isModerator = false; // Заменить на реальную проверку

    // Проверяем, является ли сообщение ответом
    const isReply = message.reply_id !== undefined && message.reply_id !== null;
    console.log(message.message_id, message.reply_id, isReply)
    let replyLevel = message.reply_level || 0; // Уровень вложенности (заглушка)

    // Ограничиваем уровень вложенности до 4
    if (replyLevel > 4) {
        replyLevel = 4;
    }

    // Добавляем классы для стилизации
    if (isReply) {
        messageElement.classList.add("reply", `level-${replyLevel}`);
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
        ${isReply ? `<div class="reply-to">в ответ на #${message.reply_id}</div>` : ''}
        <div class="message-content">${message.value}</div>
        <div class="message-timestamp">${formattedTime}</div>
        <div class="message-actions">
            <button class="reply-btn">Ответить</button>
            <button class="delete-btn">Удалить</button>
            <button class="ban-btn">Заблокировать</button>
        </div>
    `;

    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}