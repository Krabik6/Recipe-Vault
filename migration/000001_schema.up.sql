CREATE TABLE users
(
    id   BIGSERIAL primary key,
    login TEXT not null,
    password TEXT not null
);

