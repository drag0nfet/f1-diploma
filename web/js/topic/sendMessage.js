import {addMessageToDOM} from "./addMessageToDOM.js";
import {currentReplyId, resetReplyId} from "./replyBtn.js";

export function initSendMessage(topicId){
    const sendBtn = document.getElementById("send-message-btn");
    const messageInput = document.getElementById("message-text");

    sendBtn.addEventListener("click", function (event) {
        event.preventDefault();

        const content = messageInput.value.trim();

        const body = {
            chat_id: topicId,
            content: content
        };

        if (currentReplyId) {
            body.reply_id = currentReplyId;
        }
        fetch('/send-message', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify(body)
        })
            .then(response => response.json())
            .then(data => {

                if (data.success) {
                    addMessageToDOM(data.message);
                    messageInput.value = "";
                    resetReplyId()
                    const replyBanner = document.getElementById("reply-banner");
                    if (replyBanner) {
                        replyBanner.remove();
                    }
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