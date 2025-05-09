export function initDeleteBtn(event_id) {
    const delete_btn = document.querySelector(".delete-btn");

    delete_btn.innerHTML = 'Удалить ивент'
    delete_btn.addEventListener("click", function (e) {
        e.preventDefault();

        const confirmed = window.confirm("Вы уверены, что хотите удалить ивент? Отменить это действие будет невозможно");

        if (confirmed) {
            fetch(`/delete-event?event_id=${encodeURIComponent(event_id)}`, {
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
                        alert("Ивент удалён!");
                        window.location.href = `/editing_events`;
                    });
                })
                .catch(error => {
                    console.error("Ошибка:", error);
                    alert("Не удалось удалить ивент: " + error.message);
                });
        }
    });
}