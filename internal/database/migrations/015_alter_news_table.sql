DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_name = 'News'
              AND column_name = 'status'
        ) THEN
            ALTER TABLE "News"
                ADD COLUMN status VARCHAR(10) NOT NULL DEFAULT 'DRAFT'
                    CHECK (status IN ('ACTIVE', 'ARCHIVE', 'DRAFT'));
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'News'
              AND indexname = 'idx_news_created_at_status'
        ) THEN
            CREATE INDEX idx_news_created_at_status
                ON public."News" (created_at DESC, status);
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'News'
              AND indexname = 'idx_news_creator_id'
        ) THEN
            CREATE INDEX idx_news_creator_id
                ON public."News" (creator_id);
        END IF;
    END $$;