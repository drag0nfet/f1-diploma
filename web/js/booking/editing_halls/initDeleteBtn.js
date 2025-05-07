export function initDeleteBtn(hall_id) {
    const delete_btn = document.querySelector(".delete-btn");

    delete_btn.innerHTML = 'Удалить зал'
    delete_btn.addEventListener("click", function (e) {
        e.preventDefault();

        const confirmed = window.confirm("Вы уверены, что хотите удалить зал? Отменить это действие будет невозможно");

        if (confirmed) {
            fetch(`/delete-hall?hall_id=${encodeURIComponent(hall_id)}`, {                method: "DELETE",
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
                        alert("Зал удалён!");
                        window.location.href = `/editing_halls`;
                    });
                })
                .catch(error => {
                    console.error("Ошибка:", error);
                    alert("Не удалось удалить зал: " + error.message);
                });
        }
    });
}