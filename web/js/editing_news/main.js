import { initMenu }     from '../menu.js';
import {initAuthStatus} from "../checkAuth.js";
import {initButtons} from "./initButtons.js";
import {loadNewsInfo} from "./loadNewsInfo.js";

let news_id = -1;

document.addEventListener("DOMContentLoaded", async function () {
    initMenu();

    const ans = await initAuthStatus(8, "user", "news-list")

    if (!ans.isModerator) {
        alert("А вам здесь быть не положено! Авторизуйтесь!")
        window.location.href = '/';
    }

    const params = new URLSearchParams(window.location.search);
    news_id = parseInt(params.get("id"));

    loadNewsInfo(news_id)
    initButtons(news_id)
});