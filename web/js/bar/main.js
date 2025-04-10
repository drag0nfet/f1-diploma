import {initAuthStatus} from "../checkAuth.js";
import {initMenu} from "../menu.js";

let isModerator = false;

document.addEventListener("DOMContentLoaded", function () {
    // Модератор имеет в правах бит 4, нет гест-режима
    isModerator = initAuthStatus(4, "user", "bar");
    initMenu();
});