CREATE TABLE users
(
    id   BIGSERIAL primary key,
    login TEXT not null,
    password TEXT not null
);

CREATE TABLE recipes
(
  id serial primary key,
  name text not null,
  description text
);