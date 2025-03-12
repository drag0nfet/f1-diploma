export function initAuth() {
    const loginBtn = document.getElementById("login-btn");
    const usernameInput = document.getElementById("username");
    const passwordInput = document.getElementById("password");
    const greeting = document.getElementById("greeting");
    const authForm = document.getElementById("auth-form");
    const userNameDisplay = document.getElementById("user-name");

    loginBtn.addEventListener("click", function (event) {
        event.preventDefault();
        const username = usernameInput.value.trim();
        const password = passwordInput.value.trim();

        if (!username || !password) {
            alert("Пожалуйста, заполните все поля.");
            return;
        }

        fetch('/login', {
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
                    authForm.style.display = "none";
                    greeting.style.display = "block";
                    userNameDisplay.textContent = username;
                } else {
                    alert(data.message || "Ошибка авторизации.");
                }
            })
            .catch(error => {
                console.error("Login error:", error);
                alert("Ошибка при авторизации: " + error.message);
            });
    });
}