create table if not exists "Message"
(
    message_id   bigint default nextval('"Message_message_id_seq"'::regclass) not null
        primary key,
    chat_id      integer                                                      not null,
    sender_id    integer                                                      not null
        references "User",
    value        varchar(256)                                                 not null,
    message_time timestamp with time zone                                     not null,
    reply_id     integer
)
    partition by RANGE (message_id);

alter table "Message"
    owner to postgres;

create index if not exists idx_chat_id_time
    on "Message" (chat_id, message_time);

create table if not exists message_1_700000
    partition of "Message"
        (
            constraint "Message_sender_id_fkey"
                foreign key (sender_id) references "User"
            )
        FOR VALUES FROM ('1') TO ('700000');

alter table message_1_700000
    owner to postgres;

