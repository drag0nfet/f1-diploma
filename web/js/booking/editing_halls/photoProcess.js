export function initPhotoProcess() {
    const photoState = {
        existing: [],
        new: [],
    };

    document.getElementById('photo-input').addEventListener('change', function () {
        const files = Array.from(this.files);
        files.forEach(file => {
            if (!photoState.new.find(f => f.name === file.name && f.size === file.size)) {
                photoState.new.push(file);
            }
        });
        updatePhotoList();
        this.value = "";
    });

    function updatePhotoList() {
        const list = document.getElementById('photo-list');
        list.innerHTML = "";

        // Новые фотографии
        photoState.new.forEach((file, index) => {
            const item = document.createElement('li');
            item.textContent = file.name;

            const removeBtn = document.createElement('span');
            removeBtn.textContent = "✖";
            removeBtn.style.marginLeft = "10px";
            removeBtn.style.cursor = "pointer";
            removeBtn.style.color = "red";
            removeBtn.onclick = () => {
                photoState.new.splice(index, 1);
                updatePhotoList();
            };

            item.appendChild(removeBtn);
            list.appendChild(item);
        });

        // Существующие фотографии
        photoState.existing.forEach((photo, index) => {
            if (!photo.deleted) {
                const item = document.createElement('li');
                item.textContent = `Фото ${photo.id}`;

                const removeBtn = document.createElement('span');
                removeBtn.textContent = "✖";
                removeBtn.style.marginLeft = "10px";
                removeBtn.style.cursor = "pointer";
                removeBtn.style.color = "red";
                removeBtn.onclick = () => {
                    photoState.existing[index].deleted = true;
                    updatePhotoList();
                };

                item.appendChild(removeBtn);
                list.appendChild(item);
            }
        });
    }

    return photoState;
}