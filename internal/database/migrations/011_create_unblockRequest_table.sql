-- Создаём таблицу UnblockRequest
CREATE TABLE IF NOT EXISTS public."UnblockRequest" (
    request_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    message_id INTEGER,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'REJECTED', 'APPROVED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    comment TEXT,
    CONSTRAINT fk_unblock_request_forum_block_list
    FOREIGN KEY (user_id, message_id)
        REFERENCES public."ForumBlockList" (user_id, message_id)
        ON DELETE CASCADE
);

-- Устанавливаем владельца для новой таблицы
ALTER TABLE IF EXISTS public."UnblockRequest"
    OWNER TO postgres;