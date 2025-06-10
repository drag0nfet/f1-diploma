DO $$
    BEGIN
        IF EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'Item_image'
        ) THEN
            EXECUTE 'DROP TABLE public."Item_image" CASCADE';
        END IF;
    END$$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'Item'
        ) THEN
            EXECUTE 'DROP TABLE public."Item" CASCADE';
        END IF;
    END$$;

DO $$
    BEGIN
        IF EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_name = 'Purchase'
        ) THEN
            EXECUTE 'DROP TABLE public."Purchase" CASCADE';
        END IF;
    END$$;