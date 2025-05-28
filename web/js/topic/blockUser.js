import {deleteMessage} from "../deleteMessage.js";
import {showNotification} from "../notification";

export function initBlockButtons() {
    const blockButtons = document.getElementsByClassName("ban-btn");
    for (let i = 0; i < blockButtons.length; i++) {
        blockButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const messageItem = this.closest(".message-item");
            const messageId = parseInt(messageItem.querySelector(".message-id")
                .textContent.replace('#', '').trim());

            // Ограничение прав пользователя и создание записи в таблице блокировок
            blockUser(messageId, messageItem);
        });
    }
}

function blockUser(messageId, messageItem) {
    fetch(`/block-user/${messageId}`, {
        method: 'POST',
        credentials: "include",
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            return response.json().then(data => {
                if (!response.ok) {
                    throw new Error(data.message || "Неизвестная ошибка");
                }
                return data;
            });
        })
        .then(_ => {
            showNotification(
                "success",
                `Пользователь заблокирован`
            )
            deleteMessage(messageItem, messageId);
        })
        .catch(error => {
            console.error('Ошибка:', error);
            showNotification(
                "error",
                `Не удалось заблокировать пользователя: ${error.message}`
            )
        });
}