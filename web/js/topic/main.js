import { initMenu }     from '../menu.js';
import {loadTopicData}  from "./loadTopicData.js";

document.addEventListener("DOMContentLoaded", function () {
    initMenu();
    const topicId = window.location.pathname.split("/").pop(); // Извлекаем topicId из URL
    loadTopicData(topicId);
});