import {initAuthStatus} from "../checkAuth.js";
import {showNotification} from "../notification";
import {initEdit} from "../news-list/initEdit";
import {loadNews} from "../news-list/loadNews";

export function initAuth() {
    const loginBtn = document.getElementById("login-btn");
    const usernameInput = document.getElementById("username");
    const passwordInput = document.getElementById("password");
    const greeting = document.getElementById("greeting");
    const authForm = document.getElementById("auth-form");

    // Функция проверки статуса авторизации
    function checkAuthStatus() {
        fetch('/check-auth', {
            method: 'GET',
            credentials: 'include',
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        })
            .then(response => {
                if (response.status === 401) {
                    // Если не авторизован (например, при попытке доступа к /userPage), перенаправляем на /
                    window.location.href = "/";
                    return Promise.reject(new Error("Не авторизован"));
                }
                return response.json();
            })
            .then(data => {
                if (data.success && data.username) {
                    const username = data.username;
                    const ip = window.location.hostname;
                    const port = ":5051";
                    const accountLink = `http://${ip}${port}/account/${username}`;

                    // Формируем текст со ссылкой
                    greeting.innerHTML = `Привет, <a href="${accountLink}">${username}</a>!`;
                    greeting.style.display = "block";

                    // Скрываем форму авторизации
                    authForm.style.display = "none";
                } else {
                    greeting.style.display = "none";
                    authForm.style.display = "block";
                    console.log("Пользователь не авторизован:", data.message);
                }
            })
            .catch(error => {
                console.error("Ошибка проверки авторизации:", error);
                greeting.style.display = "none";
                authForm.style.display = "block";
            });
    }

    // Проверка статуса при загрузке страницы
    checkAuthStatus();

    // Обработчик логина
    loginBtn.addEventListener("click", function (event) {
        event.preventDefault();
        const username = usernameInput.value.trim();
        const password = passwordInput.value.trim();

        if (!username || !password) {
            showNotification(
                "error",
                `Заполните все поля!`
            )
            return;
        }
        let isModerator, _

        fetch('/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
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
            .then(async data => {
                if (data.success) {
                    // Обновляем статус после логина
                    checkAuthStatus();
                    let ans = await initAuthStatus(8, "guest", "index")

                    if (ans.isModerator) {
                        await loadNews("ACTIVE", 1, 10, isModerator)
                        await initEdit()
                    }
                } else {
                    showNotification(
                        "error",
                        `${data.message || "Ошибка авторизации"}`
                    )

                }
            })
            .catch(error => {
                console.error("Login error:", error);
                showNotification(
                    "error",
                    error
                )

            });
    });
}