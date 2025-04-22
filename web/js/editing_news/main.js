import { initMenu }     from '../menu.js';

let news_id = -1;

document.addEventListener("DOMContentLoaded", async function () {
    initMenu();

    const params = new URLSearchParams(window.location.search);
    news_id = parseInt(params.get("id"));
});

document.getElementById("publish-btn").addEventListener("click", function(e) {
    e.preventDefault();
    const confirmed = window.confirm("Вы уверены, что хотите опубликовать эту новость?");

    if (confirmed) {
        fetch(`/update-news-status`, {
            method: 'POST',
            credentials: "include",
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({
                news_id: news_id,
                status: "ACTIVE"
            }),
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
                alert('Новость опубликована!');
            })
            .catch(error => {
                console.error('Ошибка:', error);
                alert('Не удалось опубликовать новость: ' + error.message);
            });
    }
});

document.getElementById("draft-btn").addEventListener("click", function(e) {
    e.preventDefault();

    fetch(`/update-news-status`, {
        method: 'POST',
        credentials: "include",
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: JSON.stringify({
            news_id: news_id,
            status: "DRAFT"
        }),
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
            alert('Черновик новости сохранён');
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Не удалось сохранить черновик новости: ' + error.message);
        });
});

document.getElementById("archive-btn").addEventListener("click", function(e) {
    e.preventDefault();

    fetch(`/update-news-status`, {
        method: 'POST',
        credentials: "include",
        headers: {
            'Content-Type': 'application/json',
            'X-Requested-With': 'XMLHttpRequest'
        },
        body: JSON.stringify({
            news_id: news_id,
            status: "ARCHIVE"
        }),
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
            alert('Новость архивирована');
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Не удалось архивировать новость: ' + error.message);
        });
});

document.getElementById("delete-btn").addEventListener("click", function(e) {
    e.preventDefault();

    const confirmed = window.confirm("Вы уверены, что хотите удалить эту новость? Отменить это действие будет невозможно")

    if (confirmed) {
        fetch(`/delete-news/${news_id}`, {
            method: 'DELETE',
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
                    alert('Новость удалена!');
                    window.location.href = `/`;
                });
            })
            .catch(error => {
                console.error('Ошибка:', error);
                alert('Не удалось удалить новость: ' + error.message);
            });
    }
});