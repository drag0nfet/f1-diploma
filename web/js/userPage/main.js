import {initLogout}   from './logout.js';
import {initMenu}     from '../menu.js';
import {initAuthStatus} from "../checkAuth.js";
import {loadBlocks} from "./loadBlocks.js";

const userName = document.getElementById("user-name");
const guestContent = document.getElementById("guest-content");
const hostContent = document.getElementById("host-content");

document.addEventListener("DOMContentLoaded", async function () {
    initMenu()

    // Модератор здесь не нужен, следовательно, bit=null.
    // Если ты - хост, тебе показывается hostContent, иначе - guestContent
    let username = await initAuthStatus(null, "user", "userPage")
    userName.innerHTML = `${username}`;

    const urlUsername = window.location.pathname.split('/').pop();

    let sameUser = username === urlUsername;
    if (sameUser) {
        loadBlocks(username)
        initLogout()
        hostContent.style.display = "block";
    } else {
        guestContent.style.display = "block";
    }
});