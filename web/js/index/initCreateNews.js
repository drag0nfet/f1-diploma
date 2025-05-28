import {showNotification} from "../notification";

export function initCreateNews() {
    const title = document.getElementById("news_title").value.trim();
    if (!title) {
        showNotification(
            "error",
            `Введите название новости!`
        )

        return;
    }

    fetch('/create-news', {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: JSON.stringify({ title: title }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {showNotification(
                "success",
                `Новость создана!`
            )

                document.getElementById("news_title").value = "";
                return data.news.news_id
            } else {
                showNotification(
                    "error",
                    `Ошибка: ${data.message}`
                )
                console.error(data.message);
            }
        }).then(news_id => {
        window.location.href = `/editing_news?id=${news_id}`;
    })
        .catch(error => console.error("Ошибка создания новости:", error));
}
