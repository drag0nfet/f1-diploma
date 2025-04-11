import {deleteMessage} from "../deleteMessage.js";

export function initBlockButtons() {
    const blockButtons = document.getElementsByClassName("ban-btn");
    for (let i = 0; i < blockButtons.length; i++) {
        blockButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const messageItem = this.closest(".message-item");
            const messageId = parseInt(messageItem.querySelector(".message-id")
                .textContent.replace('#', '').trim());

            // Удаление сообщения в интерфейсе
            deleteMessage(messageItem, messageId);
            // Ограничение прав пользователя и создание записи в таблице блокировок
            blockUser(messageId);
        });
    }
}

function blockUser(messageId) {
    fetch(`/block-user/${messageId}`, {
        method: 'POST',
        credentials: "include",
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            if (!response.ok) {
                console.error('Ошибка при блокировке пользователя:', response.message)
            } else {
                alert('Вы успешно заблокировали пользователя')
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Не удалось заблокировать пользователя. Попробуйте еще раз.');
        });

}