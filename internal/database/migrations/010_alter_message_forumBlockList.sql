DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name='ForumBlockList' AND column_name='time_got'
        ) THEN
            ALTER TABLE "ForumBlockList"
                ADD COLUMN time_got timestamp with time zone DEFAULT now() NOT NULL;
        END IF;
    END
$$;
