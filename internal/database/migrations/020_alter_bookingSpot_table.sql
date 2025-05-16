DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM information_schema.table_constraints
            WHERE constraint_name = 'check_status_values'
              AND table_name = 'BookingSpot'
              AND table_schema = 'public'
        ) THEN
            ALTER TABLE public."BookingSpot"
                ADD CONSTRAINT check_status_values
                    CHECK (status IN ('ACTIVE', 'INACTIVE', 'RESERVED'));
        END IF;

        ALTER TABLE public."BookingSpot"
            ALTER COLUMN status SET DEFAULT 'INACTIVE';
    END $$;