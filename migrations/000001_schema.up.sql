CREATE TABLE users
(
    id   serial primary key,
    login TEXT not null,
    password TEXT not null
);

CREATE TABLE recipes
(
  id serial primary key,
  title text not null,
  description text
);

CREATE TABLE schedule
(
    id serial primary key,
    title text not null,
    description text
);

CREATE TABLE user_schedule
(
    id serial primary key,
    user_id int references users (id) on delete cascade,
    recipe_id int references recipes (id) on delete cascade
);

CREATE TABLE user_recipe
(
    id serial primary key,
    user_id int references users (id) on delete cascade
);
