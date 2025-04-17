import {initAuthStatus} from "../checkAuth.js";
import {initCreateTheme} from "./createTopic.js";
import {initMenu} from "../menu.js";

let isModerator = false;

document.addEventListener("DOMContentLoaded", async function () {
    // Модератор имеет в правах бит 1, есть гест-режим
    const {_isModerator, _} = await initAuthStatus(1, "guest", "forum");
    isModerator = _isModerator;
    initMenu();
});

document.getElementById("create_topic-btn").addEventListener("click", function(e) {
    e.preventDefault();
    initCreateTheme();
});