-- Table: public.Item

-- DROP TABLE IF EXISTS public."Item";

CREATE TABLE IF NOT EXISTS public."Item"
(
    item_id integer NOT NULL DEFAULT nextval('item_item_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    price numeric(10,2) NOT NULL,
    user_id integer NOT NULL,
    CONSTRAINT item_pk PRIMARY KEY (item_id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Item"
    OWNER to postgres;