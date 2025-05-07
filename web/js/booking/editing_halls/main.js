import {initAuthStatus} from "../../checkAuth";
import {initMenu} from "../../menu";
import {loadHallData} from "./loadHallData";
import {saveHall} from "./saveHall";
import {get_halls_list} from "./getHallsList";
import {initPhotoProcess} from "./photoProcess";
import {initTableActions} from "./tableActions";
import {modalTable} from "./editTable";
import {initDeleteBtn} from "./initDeleteBtn";

document.addEventListener("DOMContentLoaded", async function () {
    await initAuthStatus(16, "user", "editing_halls");
    initMenu();

    sessionStorage.setItem("hall_id", "-1");

    const photoState = initPhotoProcess();
    window.photoState = photoState;
    get_halls_list(photoState);

    document.querySelector(".create-hall-btn").addEventListener("click", function(e) {
        e.preventDefault();
        sessionStorage.setItem("hall_id", "-1");
        loadHallData(-1, photoState);

        initFooDeleteBtn()
    });

    document.querySelector(".create-table-btn").addEventListener("click", function(e) {
        e.preventDefault();
        const rect = e.target.getBoundingClientRect();
        modalTable(-1, rect.left, rect.top);
    });

    document.querySelector(".save-btn").addEventListener("click", function(e) {
        e.preventDefault();
        saveHall(photoState);

        initDeleteBtn(sessionStorage.getItem("hall_id"));
    });

    initTableActions()
});

function initFooDeleteBtn() {
    const delete_btn = document.querySelector(".delete-btn");

    delete_btn.innerHTML = 'Отменить создание'
    delete_btn.addEventListener("click", function (e) {
        e.preventDefault();
        window.location.href = `/booking`;
    });
}