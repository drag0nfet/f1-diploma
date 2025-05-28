import {showNotification} from "../notification";

export function initRegister() {
    const registerBtn = document.getElementById("register-btn");
    const usernameInput = document.getElementById("username");
    const passwordInput = document.getElementById("password");

    registerBtn.addEventListener("click", function (event) {
        event.preventDefault();
        const username = usernameInput.value.trim();
        const password = passwordInput.value.trim();

        if (!username || !password) {
            showNotification(
                "error",
                `Заполните все поля`
            )
            return;
        }

        fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })
            .then(response => {
                return response.json().then(data => {
                    if (!response.ok) {
                        throw new Error(data.message || `HTTP error! Status: ${response.status}`);
                    }
                    return data;
                });
            })
            .then(data => {
                if (data.success) {
                    alert("Регистрация прошла успешно! Теперь Вы можете авторизоваться, используя введённые ранее данные");
                    usernameInput.value = '';
                    passwordInput.value = '';
                } else {
                    showNotification(
                        "error",
                        `${data.message || "Ошибка регистрации"}`
                    )
                }
            })
            .catch(error => {
                console.error("Register error:", error);
                showNotification(
                    "error",
                    `Ошибка: ${error.message}`
                )
            });
    });
}