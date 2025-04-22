import {loadNews} from "../editing_news/loadNews.js";

const mainHeader = document.getElementById("main-head");
const newsContainer = document.getElementById("news-container");

document.addEventListener("DOMContentLoaded", function () {
    const params = new URLSearchParams(window.location.search);
    const status = params.get("status");

    if (status === "draft") {
        renderDraftNews();
    } else if (status === "archive") {
        renderArchiveNews();
    } else {
        renderDefaultNews();
    }

    newsContainer.addEventListener("click", function (e) {
        if (e.target.classList.contains("edit-btn")) {
            const newsId = e.target.getAttribute("data-news-id");
            window.location.href = '/editing_news?id=${newsId}';
        }
    });
});

function renderDraftNews() {
    mainHeader.innerHTML = "Ваши черновики новостей";
    loadNews("DRAFT", 1, 10);
}

function renderArchiveNews() {
    mainHeader.innerHTML = "Архив ваших новостей";
    loadNews("ARCHIVE", 1, 10);
}

function renderDefaultNews() {
    mainHeader.innerHTML = "А вам здесь быть не положено! Авторизуйтесь!";
}