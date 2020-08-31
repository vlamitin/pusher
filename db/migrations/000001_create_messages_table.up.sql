CREATE TABLE IF NOT EXISTS messages(
    message_id serial PRIMARY KEY,
    status smallint NOT NULL,
    send_time timestamp without time zone NOT NULL
);
