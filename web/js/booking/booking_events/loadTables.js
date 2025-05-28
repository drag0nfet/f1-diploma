import { getEventIdFromUrl, handleBooking } from "./handleBooking";

export async function loadTables(hallId) {
    const tablesGrid = document.getElementById("tables-grid");
    const selectedHall = document.getElementById("selected-hall");
    const hallName = document.getElementById("hall-name");
    const hallDescription = document.getElementById("hall-description");
    const carouselImages = document.getElementById("carousel-images");
    const prevBtn = document.getElementById("carousel-prev");
    const nextBtn = document.getElementById("carousel-next");

    try {
        // Запрос информации о зале и фотографиях
        const hallResponse = await fetch(`/booking/hall/${hallId}`, {
            method: "GET",
            headers: {
                "X-Requested-With": "XMLHttpRequest"
            }
        });
        const hallData = await hallResponse.json();
        if (hallData.success) {
            hallName.textContent = hallData.hall.name;
            hallDescription.textContent = hallData.hall.description || "Описание отсутствует";

            // Карусель фотографий
            carouselImages.innerHTML = "";
            const photos = hallData.photos || [];
            photos.forEach(photo => {
                const img = document.createElement("img");
                img.src = `data:${photo.mime_type};base64,${photo.content}`;
                img.className = "carousel-image";
                carouselImages.appendChild(img);
            });

            // Управление видимостью кнопок карусели
            if (photos.length > 3) {
                prevBtn.style.display = "block";
                nextBtn.style.display = "block";
            } else {
                prevBtn.style.display = "none";
                nextBtn.style.display = "none";
            }
        } else {
            hallName.textContent = "Ошибка загрузки зала";
            hallDescription.textContent = "";
            console.error("Ошибка загрузки зала:", hallData.message);
        }

        // Запрос столов
        const url = `/booking/hall/${hallId}/tables?event_id=${getEventIdFromUrl()}`
        const tablesResponse = await fetch(url, {
            method: "GET",
            headers: {
                "X-Requested-With": "XMLHttpRequest"
            }
        });
        const tablesData = await tablesResponse.json();
        if (tablesData.success && Array.isArray(tablesData.tables)) {
            tablesGrid.innerHTML = "";
            tablesData.tables.forEach(table => {
                const tableElement = document.createElement("div");
                tableElement.className = "table-item";
                tableElement.innerHTML = `
                    <h4 class="table-name">Стол №${table.table_name}</h4>
                    <div class="spots-container"></div>
                `;
                const spotsContainer = tableElement.querySelector(".spots-container");

                // Отображаем все места
                (table.spots || []).forEach(spot => {
                    const booking = (Array.isArray(spot.bookings) &&
                        spot.bookings.length > 0 ? spot.bookings[0] : {});
                    const isBookedByMe = booking.user_id === currentUserId();
                    const isBookedByOther = booking.user_id !== null &&
                        booking.user_id !== currentUserId() && booking.status === "ACTIVE";

                    const spotButton = document.createElement("button");
                    spotButton.className = `spot-btn ${isBookedByMe ? "booked-by-me" : isBookedByOther 
                        ? "booked-by-other" : "free"}`;
                    spotButton.textContent = `Место ${spot.spot_name}`;
                    spotButton.dataset.tableId = table.table_id;
                    spotButton.dataset.spotId = spot.spot_id;
                    if (isBookedByOther) {
                        spotButton.disabled = true;
                    } else {
                        spotButton.addEventListener("click", () =>
                            handleBooking(table.table_id, spot.spot_id, isBookedByMe));
                    }
                    spotsContainer.appendChild(spotButton);
                });

                tablesGrid.appendChild(tableElement);
            });
            selectedHall.style.display = "block";
        } else {
            tablesGrid.innerHTML = "<p>Не удалось загрузить столы.</p>";
            console.error("Ошибка загрузки столов:", tablesData.message);
        }
    } catch (error) {
        console.error("Ошибка при загрузке данных:", error);
        tablesGrid.innerHTML = "<p>Ошибка загрузки данных.</p>";
        hallName.textContent = "Ошибка загрузки зала";
        hallDescription.textContent = "";
    }

    // Логика карусели
    let currentIndex = 0;
    const images = carouselImages.querySelectorAll(".carousel-image");
    const totalImages = images.length;
    function updateCarousel() {
        const offset = -currentIndex * (100 / 3); // 3 изображения в ряду
        carouselImages.style.transform = `translateX(${offset}%)`;
    }
    prevBtn.addEventListener("click", () => {
        if (currentIndex > 0) {
            currentIndex--;
            updateCarousel();
        }
    });
    nextBtn.addEventListener("click", () => {
        if (currentIndex < totalImages - 3) {
            currentIndex++;
            updateCarousel();
        }
    });

    function currentUserId() {
        return localStorage.getItem("user_id");
    }
}