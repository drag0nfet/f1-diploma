create table if not exists "Purchase"
(
    purchase_id       integer default nextval('purchase_purchase_id_seq'::regclass) not null
        constraint purchase_pk
            primary key,
    item_id           integer                                                       not null,
    buyer_id          integer                                                       not null,
    purchase_time     timestamp with time zone,
    purchase_status   varchar(9)                                                    not null,
    status_time       timestamp                                                     not null,
    purchase_quantity integer
);

alter table "Purchase"
    owner to postgres;

