import {createDish} from "./createDish.js";
import {initMenu} from "../../menu.js";

document.addEventListener("DOMContentLoaded", function () {
    initMenu();
});

document.getElementById("create_dish-btn").addEventListener("click", function(e) {
    e.preventDefault();
    createDish();
});