import {showNotification} from "./notification";

export function initDeleteBtns() {
    const deleteButtons = document.getElementsByClassName("delete-btn");
    for (let i = 0; i < deleteButtons.length; i++) {
        deleteButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const messageItem = this.closest(".message-item");
            const messageId = parseInt(messageItem.querySelector(".message-id")
                .textContent.replace('#', '').trim());

            deleteMessage(messageItem, messageId);
        });
    }
}

export function deleteMessage(elem, messageId) {
    fetch(`/delete-message/${messageId}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            if (!response.ok) {
                console.error('Ошибка при удалении сообщения:', data)
                throw new Error('Ошибка при удалении сообщения');
            }
            elem.remove();
        })
        .catch(error => {
            console.error('Ошибка:', error);
            showNotification(
                "error",
                `Не удалось удалить сообщение. Попробуйте ещё раз`
            )

        });
}