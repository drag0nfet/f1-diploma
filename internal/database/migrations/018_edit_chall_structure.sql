-- Обновление таблицы Hall: удаляем поле album
ALTER TABLE public."Hall"
    DROP COLUMN IF EXISTS album;

-- Создание таблицы HallPhotos
CREATE TABLE IF NOT EXISTS public."HallPhotos" (
                                                   id SERIAL PRIMARY KEY,
                                                   hall_id INT NOT NULL REFERENCES public."Hall"(hall_id) ON DELETE CASCADE,
                                                   content BYTEA NOT NULL,
                                                   mime_type VARCHAR(50) NOT NULL,
                                                   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE IF EXISTS public."HallPhotos"
    OWNER TO postgres;