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

-- Version: 1.05
-- Description: Create table orders
CREATE TABLE orders (
    order_id                 UUID           NOT NULL,
    restaurant_id            UUID           NOT NULL,
    customer_name            TEXT           NOT NULL,
    customer_email           TEXT           NOT NULL,
    customer_phone           TEXT           NOT NULL,
    order_type               TEXT           NOT NULL,
    order_status             TEXT           NOT NULL,
    payment_status           TEXT           NOT NULL,
    payment_method           TEXT           NOT NULL,
    subtotal                 NUMERIC(10, 2) NOT NULL,
    delivery_fee             NUMERIC(10, 2) NOT NULL DEFAULT 0,
    tax                      NUMERIC(10, 2) NOT NULL,
    total                    NUMERIC(10, 2) NOT NULL,
    special_instructions     TEXT           NULL,
    stripe_payment_intent_id TEXT           NULL,
    date_created             TIMESTAMP      NOT NULL,
    date_updated             TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);

-- Version: 1.06
-- Description: Create table order_items
CREATE TABLE order_items (
    order_item_id        UUID           NOT NULL,
    order_id             UUID           NOT NULL,
    menu_item_id         UUID           NOT NULL,
    menu_item_name       TEXT           NOT NULL,
    menu_item_price      NUMERIC(10, 2) NOT NULL,
    quantity             INT            NOT NULL,
    special_instructions TEXT           NULL,
    date_created         TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_item_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(menu_item_id)
);

-- Version: 1.07
-- Description: Create table delivery_addresses
CREATE TABLE delivery_addresses (
    address_id            UUID      NOT NULL,
    order_id              UUID      NOT NULL UNIQUE,
    street                TEXT      NOT NULL,
    city                  TEXT      NOT NULL,
    state                 TEXT      NOT NULL,
    postal_code           TEXT      NOT NULL,
    delivery_instructions TEXT      NULL,
    date_created          TIMESTAMP NOT NULL,

    PRIMARY KEY (address_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

-- Version: 1.08
-- Description: Create table payment_transactions
CREATE TABLE payment_transactions (
    transaction_id           UUID           NOT NULL,
    order_id                 UUID           NOT NULL,
    stripe_payment_intent_id TEXT           NULL,
    stripe_charge_id         TEXT           NULL,
    payment_method           TEXT           NOT NULL,
    amount                   NUMERIC(10, 2) NOT NULL,
    currency                 TEXT           NOT NULL DEFAULT 'usd',
    status                   TEXT           NOT NULL,
    error_message            TEXT           NULL,
    metadata                 JSONB          NULL,
    date_created             TIMESTAMP      NOT NULL,
    date_updated             TIMESTAMP      NOT NULL,

    PRIMARY KEY (transaction_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

-- Version: 1.09
-- Description: Create indexes for orders
CREATE INDEX idx_orders_restaurant_id ON orders(restaurant_id);
CREATE INDEX idx_orders_customer_email ON orders(customer_email);
CREATE INDEX idx_orders_order_status ON orders(order_status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_date_created ON orders(date_created DESC);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_delivery_addresses_order_id ON delivery_addresses(order_id);
CREATE INDEX idx_payment_transactions_order_id ON payment_transactions(order_id);
CREATE INDEX idx_payment_transactions_stripe_id ON payment_transactions(stripe_payment_intent_id);
