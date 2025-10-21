-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID        NOT NULL,
	name          TEXT        NOT NULL,
	email         TEXT UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	password_hash TEXT        NOT NULL,
    department    TEXT        NULL,
    enabled       BOOLEAN     NOT NULL,
	date_created  TIMESTAMP   NOT NULL,
	date_updated  TIMESTAMP   NOT NULL,

	PRIMARY KEY (user_id)
);


-- Version: 1.02
-- Description: Create table restaurants
CREATE TABLE restaurants (
    restaurant_id UUID      NOT NULL,
    name          TEXT      NOT NULL,
    description   TEXT      NULL,
    address       TEXT      NOT NULL,
    phone         TEXT      NOT NULL,
    email         TEXT      NOT NULL,
    image_url     TEXT      NULL,
    enabled       BOOLEAN   NOT NULL,
    date_created  TIMESTAMP NOT NULL,
    date_updated  TIMESTAMP NOT NULL,

    PRIMARY KEY (restaurant_id)
);

-- Version: 1.03
-- Description: Create table categories
CREATE TABLE categories (
    category_id   UUID      NOT NULL,
    name          TEXT      NOT NULL,
    description   TEXT      NULL,
    restaurant_id UUID      NOT NULL,
    enabled       BOOLEAN   NOT NULL,
    date_created  TIMESTAMP NOT NULL,
    date_updated  TIMESTAMP NOT NULL,

    PRIMARY KEY (category_id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);

-- Version: 1.04
-- Description: Create table menu_items
CREATE TABLE menu_items (
    menu_item_id  UUID           NOT NULL,
    name          TEXT           NOT NULL,
    description   TEXT           NULL,
    price         NUMERIC(10, 2) NOT NULL,
    category_id   UUID           NOT NULL,
    restaurant_id UUID           NOT NULL,
    image_url     TEXT           NULL,
    available     BOOLEAN        NOT NULL,
    date_created  TIMESTAMP      NOT NULL,
    date_updated  TIMESTAMP      NOT NULL,

    PRIMARY KEY (menu_item_id),
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);
