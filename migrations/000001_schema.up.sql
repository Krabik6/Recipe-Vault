CREATE TABLE users
(
    "id"   serial primary key,
    "name" TEXT not null,
    "username" TEXT not null,
    "password_hash" TEXT not null
);

CREATE TABLE recipes
(
  "id" serial primary key,
  "title" text not null,
  "description" text,
  "userId" int references users (id) on delete cascade not null
);

CREATE TABLE schedule
(
    id serial primary key,
    "dateOf" date not null,
    "breakfastId" int references recipes (id) on delete cascade,
    "lunchId" int references recipes (id) on delete cascade,
    "dinnerId" int references recipes (id) on delete cascade,
    "userId" int references users (id) on delete cascade not null,
     constraint userId_dateOf_unique unique("userId", "dateOf")
);

-- alter table schedule todo way to resolve the problem apex
--     add constraint userId_dateOf_unique unique("userId", "dateOf")
