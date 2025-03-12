document.addEventListener("DOMContentLoaded", function () {
    const menuBtn = document.getElementById("menu-btn");
    const leftSidebar = document.getElementById("left-sidebar");
    const closeBtn = document.getElementById("close-btn");
    //const loginForm = document.getElementById("login-form");
    const greeting = document.getElementById("greeting");
    const authForm = document.getElementById("auth-form");
    const usernameInput = document.getElementById("username");
    const passwordInput = document.getElementById("password");
    const userNameDisplay = document.getElementById("user-name");
    const registerBtn = document.getElementById("register-btn");
    const loginBtn = document.getElementById("login-btn");

    // Открытие/закрытие бокового меню
    menuBtn.addEventListener("click", function () {
        leftSidebar.classList.toggle("active");
    });

    closeBtn.addEventListener("click", function () {
        leftSidebar.classList.remove("active");
    });

    // Регистрация
    registerBtn.addEventListener("click", function (event) {
        event.preventDefault();

        const username = usernameInput.value.trim();
        const password = passwordInput.value.trim();

        if (!username || !password) {
            alert("Пожалуйста, заполните все поля.");
            return;
        }

        fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })
            .then(response => {
                if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    usernameInput.value = '';
                    passwordInput.value = '';
                } else {
                    alert(data.message || "Ошибка регистрации.");
                }
            })
            .catch(error => {
                console.error("Register error:", error);
                alert("Ошибка при регистрации: " + error.message);
            });
    });

    // Вход
    loginBtn.addEventListener("click", function (event) {
        event.preventDefault();

        const username = usernameInput.value.trim();
        const password = passwordInput.value.trim();

        if (!username || !password) {
            console.log("Validation failed: empty fields");
            alert("Пожалуйста, заполните все поля.");
            return;
        }

        console.log("Sending login request:", { username, password });

        fetch('/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })
            .then(response => {
                console.log("Login response status:", response.status);
                if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
                return response.json();
            })
            .then(data => {
                console.log("Login response data:", data);
                if (data.success) {
                    console.log("Login successful");
                    authForm.style.display = "none";
                    greeting.style.display = "block";
                    userNameDisplay.textContent = username;
                } else {
                    console.log("Login failed:", data.message);
                    alert(data.message || "Ошибка авторизации.");
                }
            })
            .catch(error => {
                console.error("Login error:", error);
                alert("Ошибка при авторизации: " + error.message);
            });
    });
});