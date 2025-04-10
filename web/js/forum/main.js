import {initAuthStatus} from "../checkAuth.js";
import {initCreateTheme} from "./createTopic.js";
import {initMenu} from "../menu.js";

let isModerator = false;

document.addEventListener("DOMContentLoaded", function () {
    // Модератор имеет в правах бит 1, есть гест-режим
    isModerator = initAuthStatus(1, "guest", "forum");
    initMenu();
});

document.getElementById("create_topic-btn").addEventListener("click", function(e) {
    e.preventDefault();
    initCreateTheme();
});