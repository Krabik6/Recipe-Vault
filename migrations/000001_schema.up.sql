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

CREATE TABLE recipes
(
    "id" serial primary key,
    "title" text not null,
    "description" text not null,
    "user_id" int references users (id) on delete cascade not null,
    "public" bool not null default(false),
    "cost" decimal not null default(0), -- cost of the meal
    "timeToPrepare" integer default(0), -- time in minutes to prepare the meal
    "healthy" healthy_level not null default '2'
);

CREATE TABLE recipe_ingredients(
        "recipeId" int references recipes (id) on delete cascade not null,
        "ingredientId" int references ingredients (id) on delete cascade not null
);

-- // :3 todo создать табличку связывающую meal И recipes

CREATE TABLE meal
(
    "id" serial primary key,
    "name" varchar, --завтрак
    "at_time" timestamp not null, -- 10.04.2045, 10:00 по мск
    "user_id" int references users (id) on delete cascade not null, -- я (userId 5323)
    constraint userId_dateOf_unique unique("user_id", "at_time")
); -- запрос к этой таблице а еще к таблице туду :3 которая meal и recipes

-- schedule

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

-- я (userId 5323)
-- завтрак - курицу, обед - шоколадное мороженое, ночной дожер - орео
-- понедельник 10.04.2045

-- 1. пользователь хочет создать расписание на день => открывается экранчик и заполняется приемы пищи согласно шаблонам
-- 2. Пользователь изменяет приемы пищи на данный день
-- 3. Пользователь сохраняет данные

-- 1.
-- 2.
-- 3.

-- получить шаблоны


-- автозаполнять данные за пользователя

-- alter table schedule todo way to resolve the problem apex
--     add constraint userId_dateOf_unique unique("userId", "dateOf")








