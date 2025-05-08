import {loadHallData} from "./loadHallData";
import {initDeleteBtn} from "./initDeleteBtn";

export function getHallsList(photoState) {
    fetch('/get-halls-list', {
        method: 'GET',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            const hallsContainer = document.querySelector(".hall-select");

            if (data.success && Array.isArray(data.halls)) {
                // Очищаем контейнер и добавляем пустую опцию по умолчанию
                hallsContainer.innerHTML = '<option value="" selected>-- Выберите зал --</option>';

                // Добавляем залы из данных
                data.halls.forEach(hall => {
                    const option = document.createElement("option");
                    option.value = hall.hall_id;
                    option.textContent = hall.name;
                    hallsContainer.appendChild(option);
                });

                // Сбрасываем форму и скрываем её при загрузке
                document.querySelector(".selected-hall").style.display = "none";
                document.getElementById("hall-name").value = "";
                document.getElementById("hall-description").value = "";
                document.querySelector(".album-container").innerHTML = "<p>Альбом пуст</p>";
                photoState.existing = [];
                photoState.new = [];

                // Обработчик события change
                hallsContainer.addEventListener("change", function () {
                    const selectedValue = this.value;
                    if (selectedValue) {
                        document.querySelector(".selected-hall").style.display = "block";
                        loadHallData(selectedValue, photoState);
                        initDeleteBtn(sessionStorage.getItem("hall_id"));
                    } else {
                        document.querySelector(".selected-hall").style.display = "none";
                        document.getElementById("hall-name").value = "";
                        document.getElementById("hall-description").value = "";
                        document.querySelector(".album-container").innerHTML = "<p>Альбом пуст</p>";
                        photoState.existing = [];
                        photoState.new = [];
                    }
                });
            }
        })
        .catch(error => console.error("Ошибка загрузки залов:", error));
}