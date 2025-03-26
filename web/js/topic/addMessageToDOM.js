export function addMessageToDOM(message) {
    const messagesContainer = document.getElementById("messages-container");
    const messageElement = document.createElement("div");
    messageElement.className = "message-item";

    // Форматируем время в "YYYY-MM-DD HH:MM"
    const date = new Date(message.message_time);
    const formattedTime = `${date.getFullYear()}-${String(date.getMonth() + 1).
        padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).
        padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
    //console.log(message)
    messageElement.innerHTML = `
        <div class="message-author">${message.username}</div>
        <div class="message-content">${message.value}</div>
        <div class="message-timestamp">${formattedTime}</div>
    `;
    messagesContainer.appendChild(messageElement);

    // Прокручиваем вниз
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}