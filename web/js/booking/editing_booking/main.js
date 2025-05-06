import {initAuthStatus} from "../../checkAuth";
import {initMenu} from "../../menu";
import {loadHallData} from "./loadHallData";
import {saveHall} from "./saveHall";
import {get_halls_list} from "./getHallsList";
import {initPhotoProcess} from "./photoProcess";
import {initTableActions} from "./tableActions";
import {modalTable} from "./editTable";

document.addEventListener("DOMContentLoaded", async function () {
    await initAuthStatus(16, "user", "editing_booking");
    initMenu();

    sessionStorage.setItem("hall_id", "-1");

    const photoState = initPhotoProcess();
    window.photoState = photoState;
    get_halls_list(photoState);

    document.querySelector(".create-hall-btn").addEventListener("click", function(e) {
        e.preventDefault();
        loadHallData(-1, photoState);
    });

    document.querySelector(".create-table-btn").addEventListener("click", function(e) {
        e.preventDefault();
        const rect = e.target.getBoundingClientRect();
        modalTable(-1, rect.left, rect.top);
    });

    document.querySelector(".save-btn").addEventListener("click", function(e) {
        e.preventDefault();
        saveHall(photoState);
    });

    initTableActions()
});