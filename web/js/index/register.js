import {showNotification} from "../notification";

export function initRegister() {
    const registerBtn = document.getElementById("register-btn");
    const loginBtn = document.getElementById("login-btn");
    const form = document.getElementById("login-form");

    registerBtn.addEventListener("click", function (event) {
        event.preventDefault();

        // Проверяем, не была ли форма уже изменена
        if (!document.getElementById("confirm-password")) {
            // Скрываем кнопку авторизации
            loginBtn.style.display = "none";

            // Меняем текст кнопки регистрации
            registerBtn.textContent = "Подтвердить регистрацию";

            // Добавляем новые поля
            const confirmPasswordField = document.createElement("label");
            confirmPasswordField.innerHTML = '<input type="password" id="confirm-password" placeholder="Подтвердите пароль" required><br>';
            form.insertBefore(confirmPasswordField, registerBtn);

            const emailField = document.createElement("label");
            emailField.innerHTML = '<input type="email" id="email" placeholder="Электронная почта" required><br>';
            form.insertBefore(emailField, registerBtn);

            // Меняем обработчик для кнопки "Подтвердить регистрацию"
            registerBtn.removeEventListener("click", initRegister);
            registerBtn.addEventListener("click", handleRegister);
        }
    });

    function handleRegister(event) {
        event.preventDefault();
        const usernameInput = document.getElementById("username");
        const passwordInput = document.getElementById("password");
        const confirmPasswordInput = document.getElementById("confirm-password");
        const emailInput = document.getElementById("email");

        const username = usernameInput.value.trim();
        const password = passwordInput.value.trim();
        const confirmPassword = confirmPasswordInput.value.trim();
        const email = emailInput.value.trim();

        // Проверка заполнения всех полей
        if (!username || !password || !confirmPassword || !email) {
            showNotification("error", "Заполните все поля");
            return;
        }

        // Проверка совпадения паролей
        if (password !== confirmPassword) {
            showNotification("error", "Пароли не совпадают");
            return;
        }

        // Проверка формата email
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            showNotification("error", "Неверный формат электронной почты");
            return;
        }

        fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password, email })
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
                    alert("Регистрация прошла успешно! Подтвердите свою электронную почту через письмо, отправленное по вашему адресу.");
                    // Сбрасываем форму до исходного состояния
                    form.innerHTML = `
                        <label>
                            <input type="text" id="username" placeholder="Имя пользователя" required>
                        </label><br>
                        <label>
                            <input type="password" id="password" placeholder="Пароль" required>
                        </label><br>
                        <button type="submit" id="register-btn">Зарегистрироваться</button>
                        <button type="button" id="login-btn">Авторизоваться</button>
                    `;
                    // Переинициализируем обработчик
                    initRegister();
                } else {
                    showNotification("error", `${data.message || "Ошибка регистрации"}`);
                }
            })
            .catch(error => {
                console.error("Register error:", error);
                showNotification("error", `Ошибка: ${error.message}`);
            });
    }
}