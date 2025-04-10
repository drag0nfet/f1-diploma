create table if not exists "Dish"
(
    dish_id     serial
        primary key,
    name        varchar not null,
    cost        integer not null,
    description text,
    image       bytea
);

alter table "Dish"
    owner to postgres;

