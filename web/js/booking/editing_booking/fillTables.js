export function fillTables(tables, grid) {
    grid.innerHTML = "";
    tables.forEach(table => {
        const tableTile = document.createElement("div");
        if (table.seats > 10) {
            tableTile.classList.add("table-tile", "large");
        } else if (table.seats >= 5 && table.seats <= 10) {
            tableTile.classList.add("table-tile", "medium");
        } else {
            tableTile.classList.add("table-tile", "small");
        }

        tableTile.innerHTML = `
                            <h4>Стол №${table.table_name}</h4>
                            <p>Статус: ${table.price_status}</p>
                            <p>Мест: ${table.seats}</p>
                            <button class="edit-table-btn" data-table-id="${table.table_id}">Редактировать</button>
                            <button class="delete-table-btn" data-table-id="${table.table_id}">Удалить стол</button>
                        `;

        grid.appendChild(tableTile);
    });
}