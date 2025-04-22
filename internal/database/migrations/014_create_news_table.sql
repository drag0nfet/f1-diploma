CREATE TABLE IF NOT EXISTS public."News" (
    news_id SERIAL PRIMARY KEY,
    creator_id INTEGER NOT NULL,
    title VARCHAR(64) NOT NULL,
    description varchar(256),
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    image BYTEA,
    CONSTRAINT fk_news_user
        FOREIGN KEY (creator_id) REFERENCES public."User" (user_id) ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public."User"
    OWNER TO postgres;