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
            const replyId = messageIdElement.textContent.replace("#", "");

            currentReplyId = replyId;
            MakeReply(replyId);
        });
    }
}

function MakeReply(replyId) {
    let replyBanner = document.getElementById("reply-banner");

    // Если элемент не найден, создаём его заново
    if (!replyBanner) {
        console.warn("Reply banner not found, creating a new one");
        replyBanner = document.createElement("div");
        replyBanner.id = "reply-banner";
        replyBanner.className = "reply-to reply-banner";
        const messageFormContainer = document.querySelector(".message-form-container");
        if (messageFormContainer) {
            messageFormContainer.insertBefore(replyBanner, messageFormContainer.firstChild);
        } else {
            console.error("Message form container not found!");
            return;
        }
    }

    // Скрываем и очищаем перед обновлением
    replyBanner.style.display = "none";
    replyBanner.innerHTML = "";

    // Показываем и обновляем содержимое
    replyBanner.style.display = "block";
    replyBanner.innerHTML = `
        В ответ на <a href="#" class="reply-link" data-message-id="${replyId}">#${replyId}</a>
        <span class="cancel-reply">✕</span>
    `;

    const cancelBtn = replyBanner.querySelector(".cancel-reply");
    cancelBtn.addEventListener("click", () => {
        replyBanner.style.display = "none";
        replyBanner.innerHTML = "";
        currentReplyId = null;
    });

    // Привязываем обработчик для ссылки
    const replyLink = replyBanner.querySelector(".reply-link");
    replyLink.addEventListener("click", function (event) {
        event.preventDefault();
        const targetMessageId = this.getAttribute("data-message-id");
        const targetMessage = Array.from(document.querySelectorAll(".message-id"))
            .find(el => el.textContent === `#${targetMessageId}`)
            ?.closest(".message-item");
        if (targetMessage) {
            targetMessage.scrollIntoView({ behavior: "smooth", block: "center" });
        }
    });
}