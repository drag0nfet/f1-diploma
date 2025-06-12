const newsContainer = document.getElementById("news-container");

export function loadNews(status, page, limit, isModerator) {
    newsContainer.innerHTML = '<div class="spinner">Загрузка...</div>';

    const url = `/load-news-by-status?status=${encodeURIComponent(status)}&page=${page}&limit=${limit}`;
    fetch(url, {
        method: 'GET',
        credentials: "include",
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            if (response.status === 403) {
                newsContainer.innerHTML = '<p>У вас недостаточно прав.</p>';
                throw new Error('Forbidden');
            }
            return response.json()
        })
        .then(data => {
            if (data.success && Array.isArray(data.all_news)) {
                newsContainer.innerHTML = '';
                data.all_news.forEach(news => {
                    const newsElement = document.createElement("div");
                    newsElement.className = "news-item";

                    const editButtonHtml = isModerator
                        ? `<button class="edit-btn" data-news-id="${news.news_id}">Редактировать</button>`
                        : '';

                    newsElement.innerHTML = `
                        <div class="news-content" data-news-id="${news.news_id}">
                            ${news.image ? `<img src="data:image/jpeg;base64,${news.image}" alt="${news.title}" class="news-image">` : '<div class="news-no-image">Нет изображения</div>'}
                            <div class="news-info">
                                <h3 class="news-title">${news.title}</h3>
                                <p class="news-description">${news.description || news.comment.slice(0, 253) + '...' || "Описание отсутствует"}</p>
                                ${editButtonHtml}
                            </div>
                        </div>
                    `;
                    newsContainer.appendChild(newsElement);
                });

                if (data.total > 0) {
                    renderPagination(data.total, page, limit, status);
                }
            } else {
                newsContainer.innerHTML = '<p>Не удалось загрузить новости.</p>';
                console.error("Ошибка: " + data.message);
            }
        })
        .catch(error => {
            if (error.message !== 'Forbidden') {
                console.error("Ошибка загрузки новостей:", error);
                newsContainer.innerHTML = '<p>Ошибка загрузки новостей.</p>';
            }
        });
}

function renderPagination(total, currentPage, limit, status) {
    const totalPages = Math.ceil(total / limit);
    const paginationContainer = document.createElement("div");
    paginationContainer.className = "pagination";

    // Кнопка "Назад"
    const prevButton = document.createElement("button");
    prevButton.textContent = "Назад";
    prevButton.disabled = currentPage === 1;
    prevButton.addEventListener("click", () => {
        if (currentPage > 1) {
            loadNews(status, currentPage - 1, limit);
        }
    });
    paginationContainer.appendChild(prevButton);

    // Определяем диапазон страниц
    const maxVisiblePages = 9;
    let startPage, endPage;

    if (totalPages <= maxVisiblePages) {
        startPage = 1;
        endPage = totalPages;
    } else {
        const halfVisible = Math.floor(maxVisiblePages / 2);
        startPage = Math.max(1, currentPage - halfVisible);
        endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);

        // Корректируем диапазон, если он уходит за границы
        if (endPage - startPage + 1 < maxVisiblePages) {
            startPage = Math.max(1, endPage - maxVisiblePages + 1);
        }
    }

    // Кнопка для первой страницы
    if (startPage > 1) {
        const firstPageButton = document.createElement("button");
        firstPageButton.textContent = "1";
        firstPageButton.addEventListener("click", () => {
            loadNews(status, 1, limit);
        });
        paginationContainer.appendChild(firstPageButton);

        if (startPage > 2) {
            const dots = document.createElement("span");
            dots.textContent = "...";
            paginationContainer.appendChild(dots);
        }
    }

    // Кнопки для страниц в диапазоне
    for (let i = startPage; i <= endPage; i++) {
        const pageButton = document.createElement("button");
        pageButton.textContent = i;
        pageButton.className = i === currentPage ? "active" : "";
        pageButton.disabled = i === currentPage;
        pageButton.addEventListener("click", () => {
            loadNews(status, i, limit);
        });
        paginationContainer.appendChild(pageButton);
    }

    // Кнопка для последней страницы
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            const dots = document.createElement("span");
            dots.textContent = "...";
            paginationContainer.appendChild(dots);
        }

        const lastPageButton = document.createElement("button");
        lastPageButton.textContent = totalPages;
        lastPageButton.addEventListener("click", () => {
            loadNews(status, totalPages, limit);
        });
        paginationContainer.appendChild(lastPageButton);
    }

    // Кнопка "Вперёд"
    const nextButton = document.createElement("button");
    nextButton.textContent = "Вперёд";
    nextButton.disabled = currentPage === totalPages;
    nextButton.addEventListener("click", () => {
        if (currentPage < totalPages) {
            loadNews(status, currentPage + 1, limit);
        }
    });
    paginationContainer.appendChild(nextButton);

    newsContainer.appendChild(paginationContainer);
}