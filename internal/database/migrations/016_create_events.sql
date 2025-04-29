CREATE TABLE IF NOT EXISTS public."Event" (
    event_id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    time_start TIMESTAMPTZ NOT NULL,
    sport_category VARCHAR(32) NOT NULL,
    sport_type VARCHAR(32) NOT NULL,
    price_status VARCHAR(20) NOT NULL,
    duration INTEGER NOT NULL DEFAULT 90
);

ALTER TABLE IF EXISTS public."Event"
    OWNER TO postgres;


CREATE TABLE IF NOT EXISTS public."Hall" (
    hall_id SERIAL PRIMARY KEY,
    table_grid INT[][] NOT NULL,
    name VARCHAR(20) NOT NULL,
    description VARCHAR(256) NOT NULL,
    album BYTEA[],
    max_tables INT NOT NULL,
    now_tables INT DEFAULT 0 NOT NULL
);

ALTER TABLE IF EXISTS public."Hall"
    OWNER TO postgres;


CREATE TABLE IF NOT EXISTS public."Table" (
    table_id SERIAL PRIMARY KEY,
    hall_id INTEGER NOT NULL REFERENCES public."Hall" (hall_id) ON DELETE CASCADE,
    table_name INTEGER NOT NULL,
    price_status VARCHAR(20) NOT NULL
);

ALTER TABLE IF EXISTS public."Table"
    OWNER TO postgres;


CREATE TABLE IF NOT EXISTS public."Spot" (
    spot_id SERIAL PRIMARY KEY,
    table_id INTEGER NOT NULL REFERENCES public."Table" (table_id) ON DELETE CASCADE,
    spot_name INTEGER NOT NULL
);

ALTER TABLE IF EXISTS public."Spot"
    OWNER TO postgres;


CREATE TABLE IF NOT EXISTS public."BookingSpot" (
    booking_id SERIAL PRIMARY KEY,
    spot_id INTEGER NOT NULL REFERENCES public."Spot" (spot_id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES public."User" (user_id) ON DELETE SET NULL,
    event_id INTEGER REFERENCES public."Event" (event_id) ON DELETE SET NULL,
    status VARCHAR(10) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL
);

ALTER TABLE IF EXISTS public."BookingSpot"
    OWNER TO postgres;