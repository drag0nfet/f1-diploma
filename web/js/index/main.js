import { initMenu }     from '../menu/menu.js';
import { initAuth }     from './auth.js';
import { initRegister } from "./register.js";

document.addEventListener("DOMContentLoaded", function () {
    initMenu();
    initAuth();
    initRegister();
});