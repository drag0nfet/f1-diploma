export function initAuthStatus() {
    const guestContent = document.getElementById("guest-content");
    const userContent = document.getElementById("user-content");
    const usernameDisplay = document.getElementById("username-display");

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
                    guestContent.style.display = "none";
                    userContent.style.display = "block";
                    usernameDisplay.textContent = data.username;
                } else {
                    guestContent.style.display = "block";
                    userContent.style.display = "none";
                    console.log("Пользователь не авторизован:", data.message);
                }
            })
            .catch(error => {
                console.error("Ошибка проверки авторизации:", error);
                guestContent.style.display = "block";
                userContent.style.display = "none";
            });
}