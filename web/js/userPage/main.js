import {initLogout}   from './logout.js';
import {initMenu}     from '../menu.js';
import {initAuthStatus} from "../checkAuth.js";
import {loadBlocks} from "./userBlocks/loadBlocks.js";
import {loadRequests} from "./moderatorBlocks/loadRequests.js";

const guestContent = document.getElementById("guest-content");
const hostContent = document.getElementById("host-content");
const userName = document.getElementById("user-name");


document.addEventListener("DOMContentLoaded", async function () {
    initMenu()

    // Модератор форума с битом = 1.
    // Если ты - хост, тебе показывается hostContent, иначе - guestContent
    const {isModerator, username} = await initAuthStatus(1, "user", "userPage");

    const urlUsername = window.location.pathname.split('/').pop();
    userName.innerHTML = `${urlUsername}`;

    let sameUser = username === urlUsername;
    if (sameUser) {
        loadBlocks(username)
        if (isModerator) {
            loadRequests()
        }
        initLogout()
        hostContent.style.display = "block";
    } else {
        guestContent.style.display = "block";
    }
});