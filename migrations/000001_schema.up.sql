CREATE TABLE users
(
    "id"   serial primary key,
    "name" TEXT not null,
    "username" TEXT not null unique,
    "password_hash" TEXT not null
);

CREATE TABLE ingredients
(
    "id" serial primary key,
    "name" text not null,
    "description" text,
    "user_id" int references users (id) on delete cascade not null,
    "public" bool not null default(false)
);

CREATE TYPE healthy_level AS ENUM ('1', '2', '3');

create table recipes
(
    id              serial
        primary key,
    title           text                                     not null,
    description     text                                     not null,
    user_id         integer                                  not null
        references users
            on delete cascade,
    public          boolean       default false              not null,
    cost            numeric       default 0                  not null,
    "timeToPrepare" integer       default 0,
    healthy         healthy_level default '2'::healthy_level not null,
    "imageURLs"     text[]
);


CREATE TABLE recipe_images
(
    "id" serial primary key,
    "recipe_id" int references recipes (id) on delete cascade not null,
    "image_data" bytea
);



CREATE TABLE recipe_ingredients(
        "recipeId" int references recipes (id) on delete cascade not null,
        "ingredientId" int references ingredients (id) on delete cascade not null
);

CREATE TABLE meal
(
    "id" serial primary key,
    "name" varchar,
    "at_time" timestamp not null,
    "user_id" int references users (id) on delete cascade not null,
    constraint userId_dateOf_unique unique("user_id", "at_time")
);

CREATE TABLE meal_template(
      "name" varchar,
      id serial primary key,
      "time" time,
      "user_id" int references users (id) on delete cascade not null
);

CREATE TABLE mealRecipes(
    "recipeId" int references recipes (id) on delete cascade not null,
    "mealId" int references meal (id) on delete cascade not null
);
