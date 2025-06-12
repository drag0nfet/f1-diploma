import {initMenu} from "../../menu";
import {marked} from "marked";
import {showNotification} from "../../notification";

document.addEventListener("DOMContentLoaded", async function () {
    initMenu();
    loadNew();
});

function loadNew() {
    const path = window.location.pathname;
    const newsId = path.substring(path.lastIndexOf('/') + 1);

    const url = `/loadNew/${newsId}`;

    fetch(url, {
        method: 'GET',
        credentials: "include",
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка загрузки новости');
            }
            return response.json();
        })
        .then(data => {
            if (!data.Success) {
                throw new Error(data.Message || 'Ошибка загрузки новости');
            }

            const news = data.News;

            const imageElement = document.querySelector('.news-image');
            const titleElement = document.querySelector('.news-title');
            const descriptionElement = document.querySelector('.news-description');
            const contentElement = document.querySelector('.news-content');
            const dateElement = document.querySelector('.news-date');

            if (news.image) {
                imageElement.src = `data:image/jpeg;base64,${news.image}`;
            } else {
                imageElement.style.display = 'none';
            }

            titleElement.textContent = news.title;

            descriptionElement.textContent = news.description || '';

            contentElement.innerHTML = marked.parse(news.comment);

            const createdAt = new Date(news.created_at);
            dateElement.textContent = createdAt.toLocaleString('ru-RU', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        })
        .catch(error => {
            showNotification("error", "Ошибка отображения новости")
        });
}