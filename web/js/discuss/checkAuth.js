export function initAuthStatus() {
    const guestContent = document.getElementById("guest-content");
    const userContent = document.getElementById("user-content");
    const moderatorContent = document.getElementById("moderator-content"); // Добавляем
    //const usernameDisplay = document.getElementById("username-display");

    fetch('/check-auth', {
        method: 'GET',
        credentials: 'include',
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.success && data.username) {
                const rights = data.rights || 0; // По умолчанию 0, если rights отсутствует
                // usernameDisplay.textContent = data.username;

                // Скрываем гостевой контент для всех авторизованных
                guestContent.style.display = "none";

                // Показываем userContent для всех авторизованных
                userContent.style.display = "block";

                // Показываем moderatorContent только для модераторов (rights % 2 === 1)
                if (rights % 2 === 1) {
                    moderatorContent.style.display = "block";
                } else {
                    moderatorContent.style.display = "none";
                }
            } else {
                // Не авторизован
                guestContent.style.display = "block";
                userContent.style.display = "none";
                moderatorContent.style.display = "none";
                console.log("Пользователь не авторизован:", data.message);
            }
        })
        .catch(error => {
            console.error("Ошибка проверки авторизации:", error);
            guestContent.style.display = "block";
            userContent.style.display = "none";
            moderatorContent.style.display = "none";
        });
}