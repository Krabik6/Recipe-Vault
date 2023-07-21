CREATE TABLE users
(
    "id"   serial primary key,
    "name" TEXT not null,
    "username" TEXT not null unique,
    "password_hash" TEXT not null
);

create table recipes
(
    "id"              serial
        primary key,
    "title"           text                                     not null,
    "description"     text                                     not null,
    "user_id"         integer                                  not null
        references users
            on delete cascade,
    "public"          boolean       default false              not null,
    "cost"            numeric       default 0                  not null,
    "timeToPrepare" integer       default 0,
    "healthy"         integer default 0 not null,
    "imageURLs"     text[]
);

CREATE TABLE ingredients (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "price" NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    "unit" TEXT NOT NULL DEFAULT '',
    "possible_units" text[] NOT NULL DEFAULT '{}',
    "unitShort" text NOT NULL DEFAULT '',
    "unitLong" text NOT NULL DEFAULT '',
    "protein" NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    "fat" NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    "carbs" NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    "aisle" TEXT NOT NULL DEFAULT '',
    "image" TEXT NOT NULL DEFAULT '',
    "categoryPath" TEXT[] NOT NULL DEFAULT '{}',
    "consistency" TEXT NOT NULL DEFAULT '',
    "external_id" int NOT NULL DEFAULT 0
);


CREATE TABLE recipe_ingredients (
    "id" SERIAL PRIMARY KEY,
    "recipe_id" INTEGER NOT NULL REFERENCES recipes(id),
    "ingredient_id" INTEGER NOT NULL REFERENCES ingredients(id),
    "amount" NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    "unit" TEXT NOT NULL DEFAULT '',
    "price" NUMERIC(10, 2) NOT NULL DEFAULT 0.0
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
