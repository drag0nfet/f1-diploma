import {initDeleteBtn} from "./initDeleteBtn";
import {getUnMdText} from "../getUnMdText";
import {loadEventData} from "./loadEventData";

export function getEventsList() {
    fetch('/get-events-list', {
        method: 'GET',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            const eventsContainer = document.querySelector(".event-select");

            if (data.success && Array.isArray(data.events)) {
                // Очищаем контейнер и добавляем пустую опцию по умолчанию
                eventsContainer.innerHTML = '<option value="" selected>-- Выберите ивент --</option>';

                // Добавляем залы из данных
                data.events.forEach(event => {
                    const option = document.createElement("option");
                    option.value = event.event_id;
                    option.textContent = getUnMdText(event.description);
                    eventsContainer.appendChild(option);
                });

                // Сбрасываем форму и скрываем её при загрузке
                document.querySelector(".selected-event").style.display = "none";
                document.getElementById("event-description").value = "";

                // Обработчик события change
                eventsContainer.addEventListener("change", function () {
                    const selectedValue = this.value;
                    if (selectedValue) {
                        document.querySelector(".selected-event").style.display = "block";
                        loadEventData(selectedValue);
                        initDeleteBtn(sessionStorage.getItem("event_id"));
                    } else {
                        document.querySelector(".selected-event").style.display = "none";
                        document.getElementById("event-description").value = "";
                    }
                });
            }
        })
        .catch(error => console.error("Ошибка загрузки залов:", error));
}