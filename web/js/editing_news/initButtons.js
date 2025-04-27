let news_id = 0;

function loadData(formData, btn) {
    const descriptionInput = document.getElementById("news-description");
    const imageInput = document.getElementById("news-image");
    const titleInput = document.getElementById("news-title");
    const commentInput = document.getElementById("news-comment");

    formData.append("news_id", news_id);
    formData.append("title", titleInput.value.trim());
    formData.append("description", descriptionInput.value.trim());
    formData.append("comment", commentInput.value.trim());
    formData.append("status", btn.getAttribute("data-status"));

    if (imageInput.files.length > 0) {
        formData.append("image", imageInput.files[0]);
    }
}

function validateRequiredFields(is_com) {
    const titleInput = document.getElementById("news-title");
    const commentInput = document.getElementById("news-comment");

    if (!titleInput.value.trim()) {
        alert("Пожалуйста, заполните название новости.");
        return false;
    }

    if (is_com && !commentInput.value.trim()) {
        alert("Пожалуйста, заполните содержание новости.");
        return false;
    }

    return true;
}

async function handleNewsAction(statusMessage, redirectStatus) {
    if (!validateRequiredFields(redirectStatus === 'ACTIVE')) {
        return;
    }

    let new_id = -1
    const formData = new FormData();
    const btn = document.querySelector(`.action-btn[data-status="${redirectStatus}"]`);
    loadData(formData, btn);

    try {
        const response = await fetch("/update-news", {
            method: "POST",
            credentials: "include",
            headers: {
                "X-Requested-With": "XMLHttpRequest",
            },
            body: formData,
        });

        const data = await response.json();

        if (data.success) {
            new_id = data.id;
            alert(statusMessage);
            if (redirectStatus === 'ACTIVE')
                window.location.href = `/`;
            return new_id
        } else {
            alert("Ошибка: " + data.message);
        }
    } catch (error) {
        console.error("Ошибка:", error);
        alert("Произошла ошибка при обработке новости.");
    }
}

export function initButtons(newsID) {
    news_id = newsID;

    document.getElementById("publish-btn").addEventListener("click", async function () {
        await handleNewsAction("Новость успешно опубликована!", "ACTIVE");
    });

    document.getElementById("draft-btn").addEventListener("click", async function () {
        let new_id = await handleNewsAction("Новость успешно сохранена как черновик!", "DRAFT");
        if (newsID === -1) {
            window.location.href = `/editing_news?id=${new_id}`;
        }
    });

    document.getElementById("archive-btn").addEventListener("click", async function () {
        let new_id = await handleNewsAction("Новость успешно архивирована!", "ARCHIVE");
        if (newsID === -1) {
            window.location.href = `/editing_news?id=${new_id}`;
        }
    });

    const delete_btn = document.getElementById("delete-btn")

    if (newsID !== -1) {
        delete_btn.innerHTML = 'Удалить'
        delete_btn.addEventListener("click", function (e) {
            e.preventDefault();

            const confirmed = window.confirm("Вы уверены, что хотите удалить эту новость? Отменить это действие будет невозможно");

            if (confirmed) {
                fetch(`/delete-news/${newsID}`, {
                    method: "DELETE",
                    credentials: "include",
                    headers: {
                        "X-Requested-With": "XMLHttpRequest",
                    },
                })
                    .then(response => {
                        return response.json().then(data => {
                            if (!response.ok) {
                                throw new Error(data.message || "Неизвестная ошибка");
                            }
                            alert("Новость удалена!");
                            window.location.href = `/`;
                        });
                    })
                    .catch(error => {
                        console.error("Ошибка:", error);
                        alert("Не удалось удалить новость: " + error.message);
                    });
            }
        });
    } else {
        delete_btn.innerHTML = 'Отменить создание'
        delete_btn.addEventListener("click", function (e) {
            e.preventDefault();
            window.location.href = `/`;
        });
    }
}