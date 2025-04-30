import {initAuthStatus} from "../../checkAuth";
import {initMenu} from "../../menu";

document.addEventListener("DOMContentLoaded", async function () {
    // Модератор имеет в правах бит 16, нет гест-режима
    await initAuthStatus(16, "user", "editing_booking");
    initMenu();


});