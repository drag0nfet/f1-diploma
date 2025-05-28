import { loadBlocks } from "./loadBlocks.js";
import {showNotification} from "../../notification";

export function initApologizeButtons() {
    const apologizeButtons = document.getElementsByClassName("apologize-btn");
    for (let i = 0; i < apologizeButtons.length; i++) {
        apologizeButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const blockItem = this.closest(".restriction-item");
            const ghostDiv = blockItem.querySelector(".ghost");
            const [message_id, user_id] = ghostDiv.textContent.split(":").map(Number);

            initApologize(user_id, message_id);
        });
    }
}

function initApologize(user_id, message_id) {
    const modal = document.createElement("div");
    modal.className = "modal";
    modal.style.display = "flex";

    const modalContent = document.createElement("div");
    modalContent.className = "modal-content";

    const header = document.createElement("h2");
    header.className = "modal-header";
    header.textContent = "Запрос на разблокировку";
    modalContent.appendChild(header);

    const messageDiv = document.createElement("div");
    messageDiv.className = "modal-message";
    messageDiv.textContent = "Загрузка сообщения...";
    modalContent.appendChild(messageDiv);

    fetch(`/get-messages-list?ids=${message_id}`, {
        method: "GET",
        credentials: "include",
        headers: {
            "X-Requested-With": "XMLHttpRequest",
        },
    })
        .then((response) => response.json())
        .then((data) => {
            if (data.success && Array.isArray(data.messages) && data.messages.length > 0) {
                messageDiv.textContent = data.messages[0].value;
            } else {
                messageDiv.textContent = "Не удалось загрузить сообщение.";
            }
        })
        .catch((error) => {
            console.error("Ошибка загрузки сообщения:", error);
            messageDiv.textContent = "Ошибка загрузки сообщения.";
        });

    const textarea = document.createElement("textarea");
    textarea.className = "modal-textarea";
    textarea.placeholder = "Оправдательная речь";
    modalContent.appendChild(textarea);

    const submitBtn = document.createElement("button");
    submitBtn.className = "modal-submit-btn";
    submitBtn.textContent = "Отправить";
    modalContent.appendChild(submitBtn);

    modal.appendChild(modalContent);
    document.body.appendChild(modal);

    submitBtn.addEventListener("click", async () => {
        const comment = textarea.value.trim();
        if (!comment) {
            showNotification(
                "error",
                `Пожалуйста введите оправдательную речь`
            )
            return;
        }

        try {
            const response = await fetch("/submit-unblock-request", {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                    "X-Requested-With": "XMLHttpRequest",
                },
                body: JSON.stringify({ user_id, message_id, comment }),
            });
            const data = await response.json();

            if (data.success) {
                // Сначала показываем notification
                showNotification(
                    "success",
                    `Запрос на разблокировку отправлен`
                )

                // Затем обновляем блокировки
                await loadBlocks(window.location.pathname.split('/').pop());
                // Закрываем модальное окно после обновления
                modal.remove();
            } else {
                showNotification(
                    "error",
                    `Ошибка: ${data.message}`
                )
            }
        } catch (error) {
            console.error("Ошибка отправки запроса:", error);
            showNotification(
                "error",
                `Ошибка отправки запроса`
            )
        }
    });

    modal.addEventListener("click", (event) => {
        if (event.target === modal) {
            modal.remove();
        }
    });
}