export function showNotification(type, message) {
    const notifications = document.getElementById("notifications");
    const notification = document.createElement("div");
    notification.className = `notification ${type}`; // Добавляем класс типа (success или error)
    notification.textContent = message;

    // Добавляем уведомление в контейнер
    notifications.appendChild(notification);

    // Показываем уведомление
    setTimeout(() => {
        notification.classList.add("show");
    }, 100); // Небольшая задержка для анимации

    // Удаляем уведомление через 5 секунд
    setTimeout(() => {
        notification.classList.remove("show");
        setTimeout(() => {
            notification.remove();
        }, 300); // Задержка для завершения анимации исчезновения
    }, 5000);
}