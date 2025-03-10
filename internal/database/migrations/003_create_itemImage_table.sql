create table if not exists "Item_image"
(
    image_id   integer default nextval('item_image_image_id_seq'::regclass) not null
        constraint item_images_pk
            primary key,
    image_data bytea                                                        not null,
    is_primary boolean default false,
    item_id    integer                                                      not null
);

alter table "Item_image"
    owner to postgres;

