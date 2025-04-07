alter table "Message"
    add column reply_id integer;
alter table "Message"
    add CONSTRAINT "Message_reply_id_fkey" FOREIGN KEY (reply_id)
        REFERENCES public."Message" (message_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL