-- Table: public.User

-- DROP TABLE IF EXISTS public."User";

CREATE TABLE IF NOT EXISTS public."User"
(
    login text COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default" NOT NULL,
    rights integer NOT NULL DEFAULT 0,
    user_id integer NOT NULL DEFAULT nextval('"User_user_id_seq"'::regclass),
    CONSTRAINT user_pk PRIMARY KEY (user_id),
    CONSTRAINT user_unique UNIQUE (login)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."User"
    OWNER to postgres;