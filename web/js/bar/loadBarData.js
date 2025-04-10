import {loadDish} from "./loadDish.js";

export function loadBarData(isModerator) {
    fetch('/get-dishes', {
        method: 'GET',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            const dishesContainer = document.getElementById("dishes-container");
            dishesContainer.innerHTML = "";
            if (data.success && Array.isArray(data.dishes)) {
                data.dishes.forEach(dish => {
                    loadDish(dish, isModerator);
                });
            }
        })
        .catch(error => console.error("Ошибка загрузки блюд бара:", error));
}