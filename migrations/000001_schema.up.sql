CREATE TABLE users
(
    "id"   serial primary key,
    "name" TEXT not null,
    "username" TEXT not null unique,
    "password_hash" TEXT not null
);

CREATE TABLE recipes
(
  "id" serial primary key,
  "title" text not null,
  "description" text,
  "user_id" int references users (id) on delete cascade not null,
  "public" bool not null default(false)

);

CREATE TABLE schedule
(
    id serial primary key,
    "date_of" date not null,
    "breakfast_id" int references recipes (id) on delete cascade ,
    "lunch_id" int references recipes (id) on delete cascade,
    "dinner_id" int references recipes (id) on delete cascade,
    "user_id" int references users (id) on delete cascade not null,
     constraint userId_dateOf_unique unique("user_id", "date_of")
);

-- alter table schedule todo way to resolve the problem apex
--     add constraint userId_dateOf_unique unique("userId", "dateOf")








