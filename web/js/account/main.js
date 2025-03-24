import { initLogout }   from './logout.js';
import { initMenu }     from '../menu/menu.js';

document.addEventListener("DOMContentLoaded", function () {
    initLogout()
    initMenu()
});