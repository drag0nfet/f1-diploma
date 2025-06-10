import { initMenu }     from '../menu.js';
import { initAuth }     from './auth.js';
import { initRegister } from "./register.js";
import {initAuthStatus} from "../checkAuth.js";
import {loadNews} from "../news-list/loadNews.js";
import {initEdit} from "../news-list/initEdit";

document.addEventListener("DOMContentLoaded", async function () {
    initMenu();
    initAuth();
    initRegister();

    // Модератор новостей - 3 бит = 8, есть гест-режим
    let {isModerator, un} = await initAuthStatus(8, "user", "index");

    await loadNews("ACTIVE", 1, 10, isModerator)

    await initEdit()
});

document.getElementById("create_news-btn").addEventListener("click", function(e) {
    e.preventDefault();
    window.location.href = `/editing_news?id=-1`;
});

document.getElementById("draft_news-btn").addEventListener("click", function(e) {
    e.preventDefault();
    window.location.href = "/news-list?status=draft";
});

document.getElementById("archive_news-btn").addEventListener("click", function(e) {
    e.preventDefault();
    window.location.href = "/news-list?status=archive";
});