import {initAuthStatus} from "./checkAuth.js";
import {initCreateTheme} from "./createTheme.js";

let isModerator = false;

document.addEventListener("DOMContentLoaded", function () {
    isModerator = initAuthStatus();

});

document.getElementById("create_discuss-btn").addEventListener("click", function(e) {
    e.preventDefault();
    initCreateTheme(isModerator);
});