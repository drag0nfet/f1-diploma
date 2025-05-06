DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'Table'
              AND column_name = 'seats'
        ) THEN
            ALTER TABLE public."Table"
                ADD COLUMN seats INT NOT NULL DEFAULT 0;
        END IF;
    END$$;