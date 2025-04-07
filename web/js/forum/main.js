import {initAuthStatus} from "./checkAuth.js";
import {initCreateTheme} from "./createTheme.js";
import {initMenu} from "../menu/menu.js";

let isModerator = false;

document.addEventListener("DOMContentLoaded", function () {
    isModerator = initAuthStatus();
    initMenu();
});

document.getElementById("create_discuss-btn").addEventListener("click", function(e) {
    e.preventDefault();
    initCreateTheme(isModerator);
});