import { initMenu }     from '../menu/menu.js';
import {loadTopicData}  from "../discuss/loadTopicData";


document.addEventListener("DOMContentLoaded", function () {
    initMenu();
    const topicId = window.location.pathname.split("/").pop(); // Извлекаем topicId из URL
    loadTopicData(topicId);
});