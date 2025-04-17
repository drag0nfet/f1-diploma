DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name='ForumBlockList' AND column_name='status'
        ) THEN
            ALTER TABLE "ForumBlockList"
                ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'READY' CHECK (status IN ('READY', 'WAITING'));
        END IF;
    END
$$;
