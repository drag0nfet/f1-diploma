export function initBlockButtons() {
    const blockButtons = document.getElementsByClassName("ban-btn");
    for (let i = 0; i < blockButtons.length; i++) {
        blockButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const messageItem = this.closest(".message-item");
            const messageAuthorElement = messageItem.querySelector(".message-author");
            const messageAuthor = messageAuthorElement.textContent;

            blockUser(messageAuthor);
        });
    }
}

function blockUser(messageAuthor) {
    fetch(`/block-user/${messageAuthor}`, {
        method: 'POST',
        credentials: "include",
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            if (!response.ok) {
                console.error('Ошибка при блокировке пользователя:', data)
            } else {
                alert('Вы успешно заблокировали пользователя')
            }
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Не удалось заблокировать пользователя. Попробуйте еще раз.');
        });
}