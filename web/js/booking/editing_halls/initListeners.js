import {fillTables} from "./fillTables";
import {showNotification} from "../../notification";

export async function initListeners(modalOverlay, tableId) {
    const hallId = sessionStorage.getItem("hall_id");
    const tableNameInput = document.getElementById('modal-table-name');
    const priceStatusInput = document.getElementById('modal-price-status');
    const spotCountInput = document.getElementById('modal-spot-count');

    const saveBtn = document.getElementById('modal-save-btn');
    const cancelBtn = document.getElementById('modal-cancel-btn');

    // Обработчики кнопок
    saveBtn.addEventListener('click', async () => {
        if (tableNameInput.value === "" || priceStatusInput.value === "") {
            showNotification(
                "error",
                `Заполните все необходимые поля!`
            )
            return;
        }
        if (hallId === "-1") {
            showNotification(
                "error",
                `Невозможно добавить стол к несохранённому залу. Сохраните зал!`
            )
            return;
        }
        const tableData = {
            hall_id: Number(hallId),
            table_id: tableId,
            table_name: Number(tableNameInput.value),
            price_status: priceStatusInput.value,
            spot_count: Number(spotCountInput.value) || 0,
        };
        const response = await fetch(`/save-table`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
            },
            body: JSON.stringify(tableData)
        });
        const data = await response.json();
        if (data.success) {
            showNotification(
                "success",
                `Стол успешно сохранён`
            )
            const tablesGrid = document.querySelector(".tables-grid");
            fillTables(data.tables, tablesGrid);
            modalOverlay.remove();
        } else {
            showNotification(
            "error",
            `Ошибка: ${data.message}`
            )

        }
    });

    cancelBtn.addEventListener('click', () => {
        modalOverlay.remove();
    });

    modalOverlay.addEventListener('click', (e) => {
        if (e.target === modalOverlay) {
            modalOverlay.remove();
        }
    });
}
