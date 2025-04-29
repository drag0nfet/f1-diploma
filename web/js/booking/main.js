import {initAuthStatus} from "../checkAuth.js";
import {initMenu} from "../menu.js";
import flatpickr from 'flatpickr';
import { Russian } from 'flatpickr/dist/l10n/ru';
document.addEventListener("DOMContentLoaded", async function () {
    // Модератор имеет в правах бит 16, нет гест-режима
    await initAuthStatus(16, "user", "booking");
    initMenu();

    flatpickr(".calendar-input", {
        mode: "range",
        dateFormat: "d.m.Y",
        defaultDate: "today",
        locale: Russian,
        inline: true,
        showMonths: 1,
        enableTime: false,
        allowInput: false
    });
});