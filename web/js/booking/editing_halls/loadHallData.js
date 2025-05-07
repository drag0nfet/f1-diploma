import {fillTables} from "./fillTables";

export function loadHallData(hallId, photoState) {
    const nameInput = document.getElementById("hall-name");
    const descriptionInput = document.getElementById("hall-description");
    const albumContainer = document.querySelector(".album-container");
    const tablesGrid = document.querySelector(".tables-grid");

    const selectedHall = document.querySelector(".selected-hall");
    selectedHall.style.display = "block";

    sessionStorage.setItem("hall_id", hallId);

    if (hallId === -1) {
        document.querySelector(".tables-grid").innerHTML = "";
        nameInput.value = "";
        descriptionInput.value = "";
        albumContainer.innerHTML = "<p>Альбом пуст</p>";
        photoState.existing = [];
        photoState.new = [];
        return;
    }

    fetch(`/get-hall?hall_id=${encodeURIComponent(hallId)}`, {
        method: 'GET',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success && data.hall) {
                const hall = data.hall;

                nameInput.value = hall.name;
                descriptionInput.value = hall.description;

                albumContainer.innerHTML = "";
                photoState.existing = [];

                if (Array.isArray(hall.album) && hall.album.length > 0) {
                    hall.album.forEach(photo => {
                        const img = document.createElement("img");
                        img.src = `data:${photo.mime_type};base64,${photo.content}`;
                        img.alt = `Фото ${photo.id}`;
                        img.classList.add("album-image");
                        albumContainer.appendChild(img);

                        photoState.existing.push({
                            id: photo.id,
                            src: img.src,
                            deleted: false,
                        });
                    });
                } else {
                    albumContainer.innerHTML = "<p>Альбом пуст</p>";
                }

                // Отображение столов
                if (Array.isArray(hall.tables) && hall.tables.length > 0) {
                    fillTables(hall.tables, tablesGrid);
                } else {
                    tablesGrid.innerHTML = "<p>Столов нет</p>";
                }

            } else {
                console.error("Ошибка загрузки данных зала:", data.message);
            }
        })
        .catch(error => console.error("Ошибка загрузки данных зала:", error));
}