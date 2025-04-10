import {loadForumData} from "./loadForumData.js";

export function deleteTopic(event) {
    const chatId = event.target.getAttribute("data-chat-id");
    if (!chatId) {
        console.error("chat_id не найден");
        return;
    }

    if (!confirm("Вы уверены, что хотите удалить эту тему?")) {
        return;
    }

    fetch('/delete-topic', {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: JSON.stringify({ chat_id: parseInt(chatId) })
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert("Тема удалена!");
                loadForumData(true); // Перезагружаем список тем
            } else {
                alert("Ошибка: " + data.message);
            }
        })
        .catch(error => console.error("Ошибка удаления темы:", error));
}