CREATE TABLE IF NOT EXISTS products (
    id integer PRIMARY KEY AUTOINCREMENT,
    title varchar(256) NOT NULL UNIQUE,
    price decimal(10,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id integer PRIMARY KEY AUTOINCREMENT,
    name varchar(256) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS lnk_category_porduct (
    category_id integer,
    product_id integer,
    PRIMARY KEY(category_id, product_id),
    FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS roles (
    id integer PRIMARY KEY AUTOINCREMENT,
    title varchar(64) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY AUTOINCREMENT,
    username varchar(128) NOT NULL UNIQUE,
    email varchar(256) NOT NULL UNIQUE,
    password varchar(512) NOT NULL,
    role_id integer NOT NULL DEFAULT 10,
    FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE SET DEFAULT
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id integer PRIMARY KEY AUTOINCREMENT,
    user_id integer NOT NULL,
    refresh_token varchar(512) NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET DEFAULT
);

INSERT INTO products (title, price) VALUES ('IBANEZ GRG121DX-BKF', 26800.00);
INSERT INTO products (title, price) VALUES ('IBANEZ GRX70QA-TRB', 25800.00);
INSERT INTO products (title, price) VALUES ('EPIPHONE Les Paul Special Satin E1 Heritage Cherry Vintage', 44000.00);
INSERT INTO products (title, price) VALUES ('YAMAHA F310', 16990.00);

INSERT INTO categories (name) VALUES ('Электрогитары');
INSERT INTO categories (name) VALUES ('Акустические гитары');
INSERT INTO categories (name) VALUES ('IBANEZ');

INSERT INTO lnk_category_porduct (category_id,product_id) VALUES (1, 1);
INSERT INTO lnk_category_porduct (category_id,product_id) VALUES (1, 2);
INSERT INTO lnk_category_porduct (category_id,product_id) VALUES (1, 3);
INSERT INTO lnk_category_porduct (category_id,product_id) VALUES (2, 4);
INSERT INTO lnk_category_porduct (category_id,product_id) VALUES (3, 1);
INSERT INTO lnk_category_porduct (category_id,product_id) VALUES (3, 2);

INSERT INTO roles (id, title) VALUES (1, 'admin');
INSERT INTO roles (id, title) VALUES (10, 'user');

.headers on
.mode markdown
PRAGMA foreign_keys = ON; -- Включает поддержку FOREIGN KEY, без этого не работает ON DELETE CASCADE

-- Получение списка всех категорий
SELECT * FROM categories;
-- Получение списка товаров в конкретной категории (в прмере Электрогитары category_id = 1)
SELECT * FROM products as p WHERE p.id IN (SELECT product_id FROM lnk_category_porduct WHERE category_id = 1);
