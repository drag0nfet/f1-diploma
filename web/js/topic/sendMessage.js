import {addMessageToDOM} from "./addMessageToDOM.js";

export function initSendMessage(topicId){
    const sendBtn = document.getElementById("send-message-btn");
    const messageInput = document.getElementById("message-text");

    sendBtn.addEventListener("click", function (event) {
        event.preventDefault();

        const content = messageInput.value.trim();
        fetch('/send-message', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                chat_id: topicId,
                content: content
            })
        })
            .then(response => response.json())
            .then(data => {
                console.log("Server response:", data); // Выводим весь ответ
                if (data.success) {
                    addMessageToDOM(data.message);
                    messageInput.value = "";
                } else {
                    alert(data.message || "Ошибка при отправке сообщения");
                }
            })
            .catch(error => {
                console.error("Message sending error: ", error);
                alert("Ошибка при отправке сообщения: " + error.message);
            });
    });
}