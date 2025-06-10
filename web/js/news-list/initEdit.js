const newsContainer = document.getElementById("news-container");

export function initEdit() {
    newsContainer.addEventListener("click", function (e) {
        if (e.target.classList.contains("edit-btn")) {
            const newsId = e.target.getAttribute("data-news-id");
            window.location.href = `/editing_news?id=${newsId}`;
        }
    });
}