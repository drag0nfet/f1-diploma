import {loadNews} from "./loadNews.js";
import {initAuthStatus} from "../checkAuth.js";
import {initMenu} from "../menu.js";

const mainHeader = document.getElementById("main-head");
const newsContainer = document.getElementById("news-container");

function initEdit() {
    newsContainer.addEventListener("click", function (e) {
        if (e.target.classList.contains("edit-btn")) {
            const newsId = e.target.getAttribute("data-news-id");
            window.location.href = `/editing_news?id=${newsId}`;
        }
    });
}

document.addEventListener("DOMContentLoaded", function () {
    initMenu()
    const params = new URLSearchParams(window.location.search);
    const status = params.get("status");

    let {success, _} = initAuthStatus(8, "user", "news-list")
    if (!success) {
        mainHeader.innerHTML = "А вам здесь быть не положено! Авторизуйтесь!";
    }

    if (status === "draft") {
        renderDraftNews();
    } else if (status === "archive") {
        renderArchiveNews();
    } else {
        renderDefaultNews();
    }

    initEdit()
});

function renderDraftNews() {
    mainHeader.innerHTML = "Ваши черновики новостей";
    loadNews("DRAFT", 1, 10, true);
}

function renderArchiveNews() {
    mainHeader.innerHTML = "Архив ваших новостей";
    loadNews("ARCHIVE", 1, 10, true);
}

function renderDefaultNews() {
    mainHeader.innerHTML = "А вам здесь быть не положено!";
}