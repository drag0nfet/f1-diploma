import {deleteDish} from "./deleteDish.js";

export function loadDish(dish, isModerator, prepend = false) {
    const dishContainer = document.getElementById("dishes-container");
    const dishElement = document.createElement("div");
    dishElement.className = "dish-item";

    let html = `
        <div class="dish-content">
            ${dish.image ? `<img src="data:image/jpeg;base64,${dish.image}" alt="${dish.name}" class="dish-image">` : '<div class="dish-no-image">Нет изображения</div>'}
            <div class="dish-info">
                <h3 class="dish-name">${dish.name}</h3>
                <p class="dish-description">${dish.description || "Описание отсутствует"}</p>
                <p class="dish-cost">Цена: ${dish.cost} руб.</p>
            </div>
        </div>
    `;
    if (isModerator) {
        html += `<button class="delete-dish-btn" data-dish-id="${dish.dish_id}">✖</button>`;
    }
    dishElement.innerHTML = html;

    if (prepend) {
        dishContainer.prepend(dishElement);
    } else {
        dishContainer.appendChild(dishElement);
    }

    if (isModerator) {
        const deleteBtn = dishElement.querySelector(".delete-dish-btn");
        deleteBtn.addEventListener("click", deleteDish);
    }
}