import {initListeners} from "./initListeners";

export async function modalTable(tableId, clickX, clickY) {
    const hallId = sessionStorage.getItem("hall_id");
    // Подгружаем шаблон модального окна
    const response = await fetch('/modal_new_table');
    const modalHTML = await response.text();
    document.body.insertAdjacentHTML('beforeend', modalHTML);

    const modalOverlay = document.getElementById('modal-table-overlay');
    const modalContent = document.getElementById('modal-table-content');
    const modalTitle = document.getElementById('modal-title');
    const tableNameInput = document.getElementById('modal-table-name');
    const priceStatusInput = document.getElementById('modal-price-status');
    const spotCountInput = document.getElementById('modal-spot-count');

    modalOverlay.style.display = 'block';

    const contentWidth = 300;
    const contentHeight = 300;
    let left = clickX + 20;
    let top = clickY - contentHeight / 2;

    if (left + contentWidth > window.innerWidth) {
        left = window.innerWidth - contentWidth - 20;
    }
    if (top < 0) {
        top = 0;
    } else if (top + contentHeight > window.innerHeight) {
        top = window.innerHeight - contentHeight;
    }

    modalContent.style.position = 'absolute';
    modalContent.style.left = `${left}px`;
    modalContent.style.top = `${top}px`;

    let tableData = { table_name: 0, price_status: '', seats: 1, spot_count: 1 };

    // Редактирование существующего стола
    if (tableId !== -1) {
        modalTitle.textContent = 'Редактирование стола';
        const response = await fetch(`/get-hall?hall_id=${encodeURIComponent(hallId)}`, {
            method: 'GET',
            headers: { 'X-Requested-With': 'XMLHttpRequest' }
        });
        const data = await response.json();
        if (data.success && data.hall.tables) {
            const table = data.hall.tables.find(t => t.table_id === tableId);
            if (table) {
                tableData = {
                    table_name: table.table_name,
                    price_status: table.price_status,
                    seats: table.seats,
                    spot_count: await getSpotCount(tableId)
                };
            }
        }
    } else {
        // Создание нового стола
        modalTitle.textContent = 'Создание нового стола';
        const response = await fetch(`/get-hall?hall_id=${encodeURIComponent(hallId)}`, {
            method: 'GET',
            headers: { 'X-Requested-With': 'XMLHttpRequest' }
        });
        const data = await response.json();
        if (data.success && data.hall.tables) {
            const maxTableName = Math.max(0, ...data.hall.tables.map(t => t.table_name)) || 0;
            tableData.table_name = maxTableName + 1;
            tableData.spot_count = 1;
        }
    }

    // Заполнение полей
    tableNameInput.value = tableData.table_name;
    priceStatusInput.value = tableData.price_status;
    spotCountInput.value = tableData.spot_count;

    await initListeners(modalOverlay, tableId)
}

async function getSpotCount(tableId) {
    const response = await fetch(`/get-spot-count?table_id=${encodeURIComponent(tableId)}`, {
        method: 'GET',
        headers: { 'X-Requested-With': 'XMLHttpRequest' }
    });
    const data = await response.json();
    return data.success ? data.count : 0;
}
