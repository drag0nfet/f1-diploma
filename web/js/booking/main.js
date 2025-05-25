import {initAuthStatus} from "../checkAuth.js";
import {initMenu} from "../menu.js";
import flatpickr from 'flatpickr';
import { Russian } from 'flatpickr/dist/l10n/ru';
import {loadEvents} from "./loadEvents";
document.addEventListener("DOMContentLoaded", async function () {
    // Модератор имеет в правах бит 16, нет гест-режима
    await initAuthStatus(16, "user", "booking");
    initMenu();

    flatpickr(".calendar-input", {
        mode: "range",
        dateFormat: "d.m.Y",
        locale: Russian,
        inline: true,
        showMonths: 1,
        enableTime: false,
        allowInput: true,
        allowInvalidPreload: false,
        onChange: function (selectedDates, dateStr, instance) {},
        parseDate: (datestr, format) => {
            const dateRegex = /^(\d{2})\.(\d{2})\.(\d{4})$/;
            if (dateRegex.test(datestr)) {
                const [, day, month, year] = datestr.match(dateRegex);
                return new Date(year, month - 1, day);
            }
            return null;
        },
        formatDate: (date, format) => {
            const day = String(date.getDate()).padStart(2, "0");
            const month = String(date.getMonth() + 1).padStart(2, "0");
            const year = date.getFullYear();
            return `${day}.${month}.${year}`;
        }
    });
    await loadEvents(true)
});

document.getElementById("halls-btn").addEventListener("click", function(e) {
    e.preventDefault();
    window.location.href = `/editing_halls`;
});

document.getElementById("events-btn").addEventListener("click", function(e) {
    e.preventDefault();
    window.location.href = `/editing_events`;
});

document.getElementById("apply-filters-btn").addEventListener("click", function(e) {
    e.preventDefault();
    loadEvents(true)
});

document.getElementById("reset-filters-btn").addEventListener("click", function(e) {
    e.preventDefault();
    // Уничтожаем текущий экземпляр Flatpickr
    const calendarInput = document.querySelector(".calendar-input");
    const flatpickrInstance = flatpickr(calendarInput);
    flatpickrInstance.destroy();

    // Сбрасываем значение input вручную
    calendarInput.value = "";

    // Сбрасываем выбор категорий
    const categorySelect = document.querySelector(".filter-category select");
    Array.from(categorySelect.options).forEach(option => {
        option.selected = false;
    });

    // Переинициализируем Flatpickr с исходной конфигурацией
    flatpickr(".calendar-input", {
        mode: "range",
        dateFormat: "d.m.Y",
        locale: Russian,
        inline: true,
        showMonths: 1,
        enableTime: false,
        allowInput: true,
        allowInvalidPreload: false,
        onChange: function (selectedDates, dateStr, instance) {},
        parseDate: (datestr, format) => {
            const dateRegex = /^(\d{2})\.(\d{2})\.(\d{4})$/;
            if (dateRegex.test(datestr)) {
                const [, day, month, year] = datestr.match(dateRegex);
                return new Date(year, month - 1, day);
            }
            return null;
        },
        formatDate: (date, format) => {
            const day = String(date.getDate()).padStart(2, "0");
            const month = String(date.getMonth() + 1).padStart(2, "0");
            const year = date.getFullYear();
            return `${day}.${month}.${year}`;
        }
    });

    // Загружаем ивенты без фильтров
    loadEvents(true);
});