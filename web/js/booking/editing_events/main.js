import {initAuthStatus} from "../../checkAuth";
import {initMenu} from "../../menu";
import {getEventsList} from "./getEventsList";
import {loadEventData} from "./loadEventData";
import {initDeleteBtn} from "./initDeleteBtn";
import {saveEvent} from "./saveEvent";
import flatpickr from "flatpickr";
import { Russian } from 'flatpickr/dist/l10n/ru';

document.addEventListener("DOMContentLoaded", async function () {
    await initAuthStatus(16, "user", "editing_events");
    initMenu();

    sessionStorage.setItem("event_id", "-1");
    getEventsList();

    document.querySelector(".create-event-btn").addEventListener("click", function(e) {
        e.preventDefault();
        sessionStorage.setItem("event_id", "-1");
        loadEventData(-1);
        initCancelCreateBtn();
    });

    document.querySelector(".save-btn").addEventListener("click", function(e) {
        e.preventDefault();
        saveEvent();
        const eventId = sessionStorage.getItem("event_id");
        initDeleteBtn(eventId);
    });

    updateMarkdownPreview();
    toggleSportTypeField();

    document.getElementById("event-sport-category").addEventListener("change", toggleSportTypeField);

    document.getElementById("event-description").addEventListener("input", updateMarkdownPreview);

    flatpickr(".calendar-input", {
        mode: "single",
        dateFormat: "d.m.Y H:i",
        defaultDate: "today",
        locale: Russian,
        enableTime: true,
        time_24hr: true,
        inline: false,
        showMonths: 1,
        allowInput: true,
        allowInvalidPreload: false,
        parseDate: (datestr, format) => {
            const dateTimeRegex = /^(\d{2})\.(\d{2})\.(\d{4})\s(\d{2}):(\d{2})$/;
            if (dateTimeRegex.test(datestr)) {
                const [, day, month, year, hours, minutes] = datestr.match(dateTimeRegex);
                return new Date(year, month - 1, day, hours, minutes);
            }
            return null;
        },
        formatDate: (date, format) => {
            const day = String(date.getDate()).padStart(2, "0");
            const month = String(date.getMonth() + 1).padStart(2, "0");
            const year = date.getFullYear();
            const hours = String(date.getHours()).padStart(2, "0");
            const minutes = String(date.getMinutes()).padStart(2, "0");
            return `${day}.${month}.${year} ${hours}:${minutes}`;
        }
    });
});

export function updateMarkdownPreview() {
    const commentInput = document.getElementById("event-description").value;
    const previewDiv = document.getElementById("markdown-preview-content");
    previewDiv.innerHTML = marked.parse(commentInput);
}

function toggleSportTypeField() {
    const sportCategory = document.getElementById("event-sport-category").value;
    const sportTypeInput = document.getElementById("event-sport-type");
    sportTypeInput.disabled = sportCategory === "";
    if (sportCategory === "") {
        sportTypeInput.value = "";
    }
}


function initCancelCreateBtn() {
    const deleteBtn = document.querySelector(".delete-btn");

    if (!deleteBtn) {
        console.error("Кнопка удаления/отмены (.delete-btn) не найдена");
        return;
    }

    deleteBtn.innerHTML = 'Отменить создание';
    deleteBtn.addEventListener("click", function (e) {
        e.preventDefault();
        window.location.href = `/booking`;
    });
}