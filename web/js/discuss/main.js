import {initAuthStatus} from "./checkAuth.js";
import {initCreateTheme} from "./createTheme.js";

document.addEventListener("DOMContentLoaded", function () {
    initAuthStatus();
});

document.getElementById("create_discuss-btn").addEventListener("click", function(e) {
    e.preventDefault();
    initCreateTheme();
});