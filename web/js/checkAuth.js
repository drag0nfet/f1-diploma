import { loadForumData } from "./forum/loadForumData.js";
import { loadBarData } from "./bar/loadBarData.js";

const guestContent = document.getElementById("guest-content");
const userContent = document.getElementById("user-content");
const moderatorContent = document.getElementById("moderator-content");

/**
 * Отвечает за отображение 3 видов контента: guestContent, userContent, moderatorContent.
 * @param bit - отвечает за бит модератора
 * @param user - права пользователя по умолчанию - может быть "guest" или "user"
 * @param page - страница, запрашивающая изменение отображения контента.
 * @returns {Promise<{isModerator: boolean, username: string}>} - возвращает структуру с полями isModerator - наличие бита модератора, username - юзернейм пользователя
 */
export async function initAuthStatus(bit, user, page) {
    let isModerator = false;
    let data;

    try {
        const response = await fetch('/check-auth', {
            method: 'GET',
            credentials: 'include',
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        });
        data = await response.json();

        if (data.success && data.username) {
            const rights = data.rights || 0;

            if (guestContent) {
                guestContent.style.display = "none";
            }

            // Показываем userContent для всех авторизованных
            if (userContent) {
                userContent.style.display = "block";
            }

            // Показываем moderatorContent только для модераторов страницы
            isModerator = (rights & bit) === bit;

            if (moderatorContent) {
                if (isModerator) {
                    moderatorContent.style.display = "block";
                } else {
                    moderatorContent.style.display = "none";
                }
            }

            // При обращении с форума подгружаем топики форума - только для авторизованных
            if (page === "forum") {
                loadForumData(isModerator);
            }
        } else {
            // Не авторизован
            userOrGuestContent(user);
        }
    } catch (error) {
        userOrGuestContent(user);
        console.error("Ошибка авторизации:", error);
    }

    if (page === "bar") {
        loadBarData(isModerator);
    }

    return {
        isModerator: isModerator,
        username: data.username,
    };
}

function userOrGuestContent(user) {
    if (user === "guest") {
        // По умолчанию гость - показываем только гостевой контент
        if (guestContent) {
            guestContent.style.display = "block";
        }
        if (userContent) {
            userContent.style.display = "none";
        }
    } else if (user === "user") {
        // По умолчанию юзер - показываем только юзер контент
        if (guestContent) {
            guestContent.style.display = "none";
        }
        if (userContent) {
            userContent.style.display = "block";
        }
    }
    if (moderatorContent) {
        moderatorContent.style.display = "none";
    }
}