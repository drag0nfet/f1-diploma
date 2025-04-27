export function loadNewsInfo(news_id) {
    if (news_id < 0) {
        document.getElementById("news-title").value = "";
        document.getElementById("news-description").value = "";
        document.getElementById("news-comment").value = "";
        updateMarkdownPreview();
        return;
    }

    fetch(`/load-news-info?news_id=${news_id}`, {
        method: 'GET',
        credentials: "include",
        headers: {
            'X-Requested-With': 'XMLHttpRequest'
        }
    })
        .then(response => {
            if (response.status === 403) {
                alert("У вас недостаточно прав для редактирования этой новости.");
                window.location.href = '/';
                throw new Error('Forbidden');
            }
            if (response.status === 404) {
                alert("Новость не найдена.");
                window.location.href = '/';
                throw new Error('Not Found');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                const news = data.news;
                document.getElementById("news-title").value = news.title || "";
                document.getElementById("news-description").value = news.description || "";
                document.getElementById("news-comment").value = news.comment || "";
                updateMarkdownPreview();

                const previewImage = document.getElementById("current-image-preview");
                if (news.image) {
                    previewImage.src = `data:image/jpeg;base64,${news.image}`;
                    previewImage.style.display = "block";
                } else {
                    previewImage.style.display = "none";
                }
            } else {
                alert("Ошибка загрузки новости: " + data.message);
                window.location.href = '/';
            }
        })
        .catch(error => {
            if (error.message !== 'Forbidden' && error.message !== 'Not Found') {
                console.error("Ошибка загрузки новости:", error);
                alert("Произошла ошибка при загрузке новости.");
                window.location.href = '/';
            }
        });
}