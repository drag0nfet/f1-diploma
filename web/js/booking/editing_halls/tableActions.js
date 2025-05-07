import {modalTable} from "./editTable";

export function initTableActions() {
    const tablesGrid = document.querySelector(".tables-grid");

    // Обработчик для кнопки "Редактировать"
    tablesGrid.addEventListener("click", function (e) {
        if (e.target.classList.contains("edit-table-btn")) {
            const tableId = e.target.getAttribute("data-table-id");
            const rect = e.target.getBoundingClientRect();
            modalTable(parseInt(tableId), rect.left, rect.top);
        }
    });

    // Обработчик для кнопки "Удалить стол"
    tablesGrid.addEventListener("click", function (e) {
        if (e.target.classList.contains("delete-table-btn")) {
            const tableId = e.target.getAttribute("data-table-id");
            if (confirm(`Вы уверены, что хотите удалить стол с ID ${tableId}?`)) {
                fetch(`/delete-table?table_id=${encodeURIComponent(tableId)}`, {
                    method: 'DELETE',
                    headers: {
                        'X-Requested-With': 'XMLHttpRequest'
                    }
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            alert("Стол успешно удалён!");
                            const hallId = sessionStorage.getItem("hall_id");
                            import("./loadHallData.js").then(module => {
                                module.loadHallData(hallId, window.photoState);
                            });
                        } else {
                            alert("Ошибка: " + data.message);
                        }
                    })
                    .catch(error => console.error("Ошибка удаления стола:", error));
            }
        }
    });
}