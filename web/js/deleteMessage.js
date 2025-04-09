export function initDeleteBtns() {
    const deleteButtons = document.getElementsByClassName("delete-btn");
    for (let i = 0; i < deleteButtons.length; i++) {
        deleteButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const messageElement = this.closest(".message-item");
            const messageIdElement = messageElement.querySelector(".message-id");

            deleteMessage(messageElement, messageIdElement);
        });
    }
}

export function deleteMessage(elem, id) {
    const messageId = parseInt(id.textContent.replace('#', ''));

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
            alert('Не удалось удалить сообщение. Попробуйте еще раз.');
        });
}