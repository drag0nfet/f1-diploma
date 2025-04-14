import { initApologiseButtons } from "./initApologiseButtons.js";

const blocksContainer = document.getElementById("restrictions-container");
const blocksPart = document.querySelector(".restriction-block");
const header = document.querySelector(".restrictions-title");
let messageIdMs = [];
let blockMs = [];

export function loadBlocks(username) {
    blocksPart.style.display = "block";
    fetch(`/get-blocks/${username}`, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            blocksContainer.innerHTML = "";
            messageIdMs = [];
            blockMs = [];

            if (data.success && Array.isArray(data.blocks) && data.blocks.length > 0) {
                data.blocks.forEach(block => {
                    if (block.is_valid) {
                        messageIdMs.push(block.message_id);
                        blockMs.push(block);
                    }
                });
                if (messageIdMs.length > 0) {
                    showBlocks();
                } else {
                    header.innerHTML = "<p>Нет активных ограничений.</p>";
                }
            } else {
                header.innerHTML = "<p>Нет активных ограничений.</p>";
            }
        })
        .catch(error => {
            console.error("Ошибка загрузки блокировок:", error);
            header.innerHTML = "<p>Не удалось загрузить ограничения.</p>";
        });
}

function showBlocks() {
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
                blockMs.forEach((block, index) => {
                    const message = data.messages.find(msg => msg.message_id === block.message_id);
                    if (message) {
                        showBlock(message, block, index + 1);
                    }
                });
                initApologiseButtons();
            } else {
                header.innerHTML = "<p>Не удалось загрузить сообщения для ограничений.</p>";
            }
        })
        .catch(error => {
            console.error("Ошибка загрузки сообщений:", error);
            header.innerHTML = "<p>Не удалось загрузить сообщения для ограничений.</p>";
        });
}

function showBlock(message, block, violationNumber) {
    const blockElement = document.createElement("div");
    blockElement.className = "restriction-item";

    const formattedTime = new Date(block.time_got).toLocaleString();

    blockElement.innerHTML = `
        <div class="restriction-header">
            <span class="violation-number">Нарушение #${violationNumber}</span>
        </div>
        <div class="restriction-message">Сообщение: ${message.value}</div>
        <div class="restriction-timestamp">Время блокировки: ${formattedTime}</div>
        <div class="restriction-actions">
            <button class="apologise-btn" data-block-id="${block.message_id}">Обжаловать</button>
        </div>
    `;

    blocksContainer.appendChild(blockElement);
}