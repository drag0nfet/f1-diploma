DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name='Message' AND column_name='is_deleted'
        ) THEN
            ALTER TABLE "Message"
                ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;
        END IF;
    END
$$;
