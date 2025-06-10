DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'User'
              AND column_name = 'email'
        ) THEN
            ALTER TABLE public."User"
                ADD COLUMN email VARCHAR(50) NOT NULL DEFAULT '';
        END IF;
    END$$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'User'
              AND column_name = 'is_confirmed'
        ) THEN
            ALTER TABLE public."User"
                ADD COLUMN is_confirmed BOOLEAN NOT NULL DEFAULT FALSE;
        END IF;
    END$$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'User'
              AND column_name = 'confirmation_token'
        ) THEN
            ALTER TABLE public."User"
                ADD COLUMN confirmation_token VARCHAR(20);
        END IF;
    END$$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.columns
            WHERE table_schema = 'public'
              AND table_name = 'User'
              AND column_name = 'last_sent'
        ) THEN
            ALTER TABLE public."User"
                ADD COLUMN last_sent TIMESTAMPTZ DEFAULT NOW();
        END IF;
    END$$;