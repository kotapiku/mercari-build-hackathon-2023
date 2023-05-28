CREATE TABLE IF NOT EXISTS category
(
    id   integer primary key,
    name varchar(50)
);

CREATE TABLE IF NOT EXISTS items
(
    id          integer primary key autoincrement,
    name        varchar(50),
    price       integer,
    description text,
    category_id integer,
    seller_id   integer,
    image       blob,
    status      integer,
    created_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    updated_at  text NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    FOREIGN KEY(category_id) REFERENCES category(id)
);

CREATE TABLE IF NOT EXISTS users
(
    id       integer primary key autoincrement,
    name     varchar(50) unique,
    password binary(60),
    balance  integer default 0
);

CREATE TABLE IF NOT EXISTS status
(
    id   integer primary key,
    name varchar(50)
);