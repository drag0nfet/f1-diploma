export let currentReplyId = null;

// Сброс переменной для ответа
export function resetReplyId() {
    currentReplyId = null;
}

// Изменение значения переменной для ответа
export function setReplyId(id) {
    currentReplyId = id;
}

export function initReplyBtn() {
    const replyButtons = document.getElementsByClassName("reply-btn");

    for (let i = 0; i < replyButtons.length; i++) {
        replyButtons[i].addEventListener("click", function (event) {
            event.preventDefault();

            const messageElement = this.closest(".message-item");
            const messageIdElement = messageElement.querySelector(".message-id");
            const replyId = messageIdElement.textContent.replace("#", ""); // Извлекаем ID (например, "123" из "#123")

            currentReplyId = replyId;

            MakeReply(replyId);
        });
    }
}

function MakeReply(replyId) {
    const existingBanner = document.getElementById("reply-banner");
    if (existingBanner) {
        existingBanner.remove();
    }

    const replyBanner = document.createElement("div");
    replyBanner.id = "reply-banner";
    replyBanner.className = "reply-banner";
    replyBanner.innerHTML = `
        <span>в ответ на #${replyId}</span>
        <span class="cancel-reply">✕</span>
    `;

    const messageForm = document.getElementById("message-form");
    messageForm.parentNode.insertBefore(replyBanner, messageForm);

    const cancelBtn = replyBanner.querySelector(".cancel-reply");
    cancelBtn.addEventListener("click", () => {
        replyBanner.remove();
        currentReplyId = null;
    });
}