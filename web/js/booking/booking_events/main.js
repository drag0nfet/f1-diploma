import { initAuthStatus } from "../../checkAuth.js";
import { initMenu } from "../../menu.js";
import {loadTables} from "./loadTables";
import {loadHalls} from "./loadHalls";

const hallSelect = document.getElementById("hall-select");
const selectedHall = document.getElementById("selected-hall");

document.addEventListener("DOMContentLoaded", async function () {
    await initAuthStatus(-1, "user", "booking_events");
    initMenu();

    await loadHalls();

    hallSelect.addEventListener("change", async function () {
        const hallId = hallSelect.value;
        if (hallId) {
            await loadTables(hallId);
        } else {
            selectedHall.style.display = "none";
        }
    });
});

