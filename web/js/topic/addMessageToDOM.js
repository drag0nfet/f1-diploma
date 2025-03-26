export function addMessageToDOM(message) {
    const messagesContainer = document.getElementById("messages-container");
    const messageElement = document.createElement("div");
    messageElement.className = "message-item";
    messageElement.innerHTML = `
        <div class="message-author">User ${message.sender_id}</div>
        <div class="message-content">${message.value}</div>
        <div class="message-timestamp">${new Date(message.message_time).toLocaleString()}</div>
    `;
    messagesContainer.appendChild(messageElement);

    // Прокручиваем вниз
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}