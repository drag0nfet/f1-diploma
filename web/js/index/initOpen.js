const newsContainer = document.getElementById("news-container");

export function initOpen() {
    newsContainer.addEventListener("click", function (e) {
        const target = e.target.closest('.news-content'); // Ищем ближайший элемент с классом news-content
        if (target && !e.target.classList.contains("edit-btn")) {
            const newsId = target.getAttribute("data-news-id");
            if (newsId) {
                window.location.href = `/news/${newsId}`;
            }
        }
    });
}

