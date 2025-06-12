import {loadNews} from "./loadNews.js";
import {initAuthStatus} from "../checkAuth.js";
import {initMenu} from "../menu.js";
import {initEdit} from "./initEdit";

const mainHeader = document.getElementById("main-head");
document.addEventListener("DOMContentLoaded", function () {
    initMenu()
    const params = new URLSearchParams(window.location.search);
    const status = params.get("status");

    let ans = initAuthStatus(8, "user", "news-list")

    if (!ans.isModerator) {
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