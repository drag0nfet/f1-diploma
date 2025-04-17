UPDATE public."UnblockRequest"
SET message_id = -1
WHERE message_id IS NULL;

-- Проверяем, есть ли NULL значения в comment, и устанавливаем пустую строку, если нужно
UPDATE public."UnblockRequest"
SET comment = ''
WHERE comment IS NULL;

-- Добавляем ограничение NOT NULL для message_id
ALTER TABLE public."UnblockRequest"
    ALTER COLUMN message_id SET NOT NULL;

-- Добавляем ограничение NOT NULL для comment
ALTER TABLE public."UnblockRequest"
    ALTER COLUMN comment SET NOT NULL;