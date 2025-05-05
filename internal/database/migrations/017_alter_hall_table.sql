ALTER TABLE public."Hall"
    DROP COLUMN IF EXISTS max_tables,
    DROP COLUMN IF EXISTS now_tables,
    DROP COLUMN IF EXISTS table_grid;
