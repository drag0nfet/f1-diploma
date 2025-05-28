export async function loadHalls() {
    const hallSelect = document.getElementById("hall-select");

    try {
        const response = await fetch("/get-halls-list", {
            method: "GET",
            headers: {
                "X-Requested-With": "XMLHttpRequest"
            }
        });
        const data = await response.json();
        if (data.success && Array.isArray(data.halls)) {
            hallSelect.innerHTML = '<option value="">-- Выберите зал --</option>';
            data.halls.forEach(hall => {
                const option = document.createElement("option");
                option.value = hall.hall_id;
                option.textContent = hall.name;
                hallSelect.appendChild(option);
            });
        } else {
            console.error("Ошибка загрузки залов:", data.message);
        }
    } catch (error) {
        console.error("Ошибка при загрузке залов:", error);
    }
}