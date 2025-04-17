import {initAuthStatus} from "../checkAuth.js";
import {initMenu} from "../menu.js";

document.addEventListener("DOMContentLoaded", async function () {
    // Модератор имеет в правах бит 4, нет гест-режима
    await initAuthStatus(4, "user", "bar");
    initMenu();
});