create table if not exists "Chat"
(
    chat_id    serial
        primary key,
    chat_type  varchar(20)                            not null,
    item_id    integer
                                                      references "Item"
                                                          on delete set null,
    created_at timestamp with time zone default now() not null,
    title      varchar(100)
);

alter table "Chat"
    owner to postgres;

create index if not exists idx_chat_type
    on "Chat" (chat_type);

