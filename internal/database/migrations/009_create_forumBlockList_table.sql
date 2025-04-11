-- Table: public.ForumBlockList

-- DROP TABLE IF EXISTS public."ForumBlockList";

CREATE TABLE IF NOT EXISTS public."ForumBlockList"
(
    user_id integer NOT NULL REFERENCES public."User" (user_id) ON DELETE CASCADE,
    message_id integer NOT NULL REFERENCES public."Message" (message_id) ON DELETE SET NULL,
    moderator_id integer NOT NULL REFERENCES public."User" (user_id) ON DELETE SET NULL,
    is_valid bool DEFAULT TRUE,
    CONSTRAINT ForumBlockList_pk PRIMARY KEY (user_id, message_id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."ForumBlockList"
    OWNER to postgres;