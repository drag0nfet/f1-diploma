export function initLogout(){
    const logoutBtn = document.getElementById("logout-btn");

    logoutBtn.addEventListener("click", function (event) {
        event.preventDefault();

        fetch('/logout', {
            method: 'POST',
            credentials: 'include' // Отправляем куки с запросом
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Успешный разлогин, перенаправляем на главную
                    window.location.href = "/";
                } else {
                    alert(data.message || "Ошибка при выходе из системы.");
                }
            })
            .catch(error => {
                console.error("Logout error:", error);
                alert("Ошибка при выходе из системы: " + error.message);
            });
    });
}