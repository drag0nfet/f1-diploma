import { loadMessage } from "./loadMessage.js";
import { currentReplyId, resetReplyId } from "./reply.js";
import { initReplyBtn } from "./reply.js";
import {showNotification} from "../notification";

export function initSendMessage(topicId, rights) {
    const sendBtn = document.getElementById("send-message-btn");
    const messageInput = document.getElementById("message-text");

    sendBtn.addEventListener("click", function (event) {
        event.preventDefault();
        if (rights % 4 / 2 === 0) {
            showNotification(
                "error",
                `Вы заблокированы за нарушение правил. Обратитесь к администратору`
            )
            return;
        }
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
                    loadMessage(data.message, rights);
                    messageInput.value = "";
                    resetReplyId();
                    const replyBanner = document.getElementById("reply-banner");
                    if (replyBanner) {
                        replyBanner.style.display = "none";
                        replyBanner.innerHTML = "";
                    }
                    initReplyBtn();
                } else {
                    showNotification(
                        "error",
                        `Ошибка: ${data.message}`
                    )

                }
            })
            .catch(error => {
                console.error("Message sending error: ", error);
                showNotification(
                    "error",
                    `Ошибка: ${error.message}`
                )
            });
    });
}