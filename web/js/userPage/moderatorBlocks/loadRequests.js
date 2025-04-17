import {initButtons} from "./initButtons.js";

const requestsContainer = document.getElementById("requests-container");
const requestsBlock = document.querySelector(".requests-block");
const header = document.querySelector(".requests-title");
let messageIdMs = [];
let requestsMs = [];

export function loadRequests() {
    requestsBlock.style.display = "none";
    fetch(`/get-requests`, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            requestsContainer.innerHTML = "";
            messageIdMs = [];
            requestsMs = [];

            if (data.success && Array.isArray(data.requests) && data.requests.length > 0) {
                data.requests.forEach(block => {
                    if (block.status === "PENDING") {
                        messageIdMs.push(block.message_id);
                        requestsMs.push(block);
                    }
                });
                if (messageIdMs.length > 0) {
                    showRequests();
                }
            }
        })
        .catch(error => {
            console.error("Ошибка загрузки блокировок:", error);
        });
}

function showRequests() {
    requestsBlock.style.display = "block";
    const idsParam = messageIdMs.join(",");
    fetch(`/get-messages-list?ids=${idsParam}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success && Array.isArray(data.messages)) {
                requestsMs.forEach((request) => {
                    const message = data.messages.find(msg => msg.message_id === request.message_id);
                    if (request.status === "PENDING" && message) {
                        showRequest(message, request);
                    }
                });
            } else {
                header.innerHTML = "<p>Не удалось загрузить сообщения для запросов на разблокировку.</p>";
            }
        }).then(_ => {
        initButtons();
    })
        .catch(error => {
            console.error("Ошибка загрузки сообщений:", error);
            header.innerHTML = "<p>Не удалось загрузить сообщения для запросов на разблокировку.</p>";
        });
}

function showRequest(message, request, ) {
    const requestElement = document.createElement("div");
    requestElement.className = "restriction-item";

    const formattedTime = new Date(request.created_at).toLocaleString();

    requestElement.innerHTML = `
        <div class="restriction-header">
            <span class="violation-number">Нарушение пользователя#${request.user_id}</span>
        </div>
        <div class="block-message">Заблокированное сообщение: ${message.value}</div>
        <div class="request-message">Сообщение реквеста: ${request.comment}</div>
        <div class="restriction-timestamp">Время запроса на разблокировку: ${formattedTime}</div>
        <div class="restriction-actions">
            <button class="approve-btn" 
                    data-request-id="${request.request_id}" 
                    data-user-id="${request.user_id}" 
                    data-message-id="${request.message_id}">Принять апелляцию</button>
            <button class="reject-btn" 
                    data-request-id="${request.request_id}" 
                    data-user-id="${request.user_id}" 
                    data-message-id="${request.message_id}">Отклонить апелляцию</button>
        </div>
    `;

    requestsContainer.appendChild(requestElement);
}