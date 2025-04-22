export function createDish() {
    const form = document.getElementById("dish-form");
    if (!form) {
        console.error("Форма не найдена!");
        return;
    }

    const formData = new FormData(form);

    fetch("/create-dish", {
        method: "POST",
        body: formData,
        headers: {
            "X-Requested-With": "XMLHttpRequest"
        }
    })
        .then(response => {
            return response.json().then(data => {
                if (!response.ok) {
                    throw new Error(data.message || "Ошибка при добавлении блюда");
                }
                alert("Блюдо успешно добавлено!\nПереадресация на страницу бара...");
                window.location.href = "web/pages/bar";
            });
        })
        .catch(error => {
            console.error("Ошибка:", error);
            alert(`Не удалось добавить блюдо: ${error.message}`);
        });
}