import {loadRequests} from "./loadRequests.js";
import {showNotification} from "../../notification";

export function initButtons() {
    const approveButtons = document.getElementsByClassName("approve-btn");
    for (let i = 0; i < approveButtons.length; i++) {
        approveButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const user_id = Number(this.dataset.userId);
            const message_id = Number(this.dataset.messageId);
            const request_id = Number(this.dataset.requestId);

            approve(user_id, message_id, request_id);
        });
    }

    const rejectButtons = document.getElementsByClassName("reject-btn");
    for (let i = 0; i < rejectButtons.length; i++) {
        rejectButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const user_id = Number(this.dataset.userId);
            const message_id = Number(this.dataset.messageId);
            const request_id = Number(this.dataset.requestId);

            reject(user_id, message_id, request_id);
        });
    }
}

function approve(user_id, message_id, request_id) {
    const body = JSON.stringify({
        user_id: user_id,
        message_id: message_id,
    });
    fetch(`/approve/${request_id}`, {
        method: "POST",
        credentials: "include",
        headers: {
            "X-Requested-With": "XMLHttpRequest",
        },
        body: body,
    })
        .then((response) => response.json())
        .then((data) => {
            if (data.success) {
                showNotification(
                    "success",
                    `Блокировка успешно снята`
                )
            } else {
                showNotification(
                    "error",
                    `Ошибка при снятии блокировки`
                )
                console.error("Ошибка снятия блокировки, бэк:", data.message);
            }
        }).then(_ => {
            loadRequests()
    })
        .catch((error) => {
            console.error("Ошибка снятия блокировки, фронт:", error);
            showNotification(
                "error",
                `Ошибка при снятии блокировки`
            )
        });
}

function reject(user_id, message_id, request_id) {
    const body = {
        user_id: user_id,
        message_id: message_id,
    }
    fetch(`/reject/${request_id}`, {
        method: "POST",
        credentials: "include",
        headers: {
            "X-Requested-With": "XMLHttpRequest",
        },
        body: JSON.stringify(body),
    })
        .then((response) => response.json())
        .then((data) => {
            if (data.success) {
                showNotification(
                    "success",
                    `Запрос на разблокировку отклонён`
                )
            } else {
                showNotification(
                    "error",
                    `Ошибка при отклонении запроса на разблокировку`
                )

                console.error("Ошибка при отклонении запроса на разблокировку, бэк:", data.message);
            }
        }).then(_ => {
        loadRequests()
    })
        .catch((error) => {
            console.error("Ошибка при отклонении запроса на разблокировку, фронт:", error);
            showNotification(
                "error",
                `Ошибка при отклонении запроса на разблокировку`
            )
        });
}