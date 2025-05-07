export async function saveHall(photoState) {
    const hallId = parseInt(sessionStorage.getItem("hall_id"));
    const name = document.getElementById('hall-name').value.trim();
    const description = document.getElementById('hall-description').value.trim();

    const formData = new FormData();
    formData.append("hall_id", hallId);
    formData.append("name", name);
    formData.append("description", description);

    // Добавляем новые фотографии
    photoState.new.forEach(file => {
        formData.append("photos", file);
    });

    // Добавляем ID удалённых фотографий
    const deletedPhotoIds = photoState.existing
        .filter(photo => photo.deleted)
        .map(photo => photo.id);
    formData.append("deleted_photo_ids", JSON.stringify(deletedPhotoIds));

    try {
        const response = await fetch('/save-hall', {
            method: 'POST',
            headers: {
                'X-Requested-With': 'XMLHttpRequest',
            },
            body: formData
        });

        const result = await response.json();
        if (result.success) {
            alert('Зал успешно сохранён!');
            // Задаём новый сессионный id зала
            sessionStorage.setItem("hall_id", result.hall.hall_id);
            // Обновляем photoState
            photoState.new = [];
            photoState.existing = result.hall.album.map(photo => ({
                id: photo.id,
                src: `data:${photo.mime_type};base64,${photo.content}`,
                deleted: false,
            }));
            // Обновляем albumContainer
            const albumContainer = document.querySelector(".album-container");
            albumContainer.innerHTML = result.hall.album.length > 0
                ? result.hall.album.map(photo => `<img src="data:${photo.mime_type};base64,${photo.content}" alt="Фото ${photo.id}" class="album-image">`).join('')
                : "<p>Альбом пуст</p>";
        } else {
            alert('Ошибка: ' + result.message);
        }
    } catch (err) {
        console.error('Ошибка при сохранении зала:', err);
        alert('Ошибка при сохранении зала');
    }
}