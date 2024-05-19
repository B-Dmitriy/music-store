CREATE TABLE IF NOT EXIST products (
    id integer primary key,
    title varchar(256) NOT NULL UNIQUE,
    price decimal(10,2) NOT NULL
);

CREATE TABLE IF NOT EXIST category (
    id integer primary key,
    name varchar(256) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXIST lnk_category_porduct (
    category_id integer,
    product_id integer,
    UNIQUE (category_id, product_id)
);

CREATE TABLE IF NOT EXIST users (
    id integer primary key,
    username varchar(128) NOT NULL UNIQUE,
    email varchar(256) NOT NULL UNIQUE,
    password varchar(512) NOT NULL,
    role_id integer
);

CREATE TABLE IF NOT EXIST roles (
    id integer primary key,
    title varchar(64) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXIST refresh_tokens (
    id integer primary key,
    user_id integer NOT NULL,
    refresh_token varchar(512) NOT NULL,
);