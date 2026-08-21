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

-- Version: 1.10
-- Description: Create table addons for category addons
CREATE TABLE addons (
    addon_id      UUID           NOT NULL,
    category_id   UUID           NOT NULL,
    restaurant_id UUID           NOT NULL,
    name          TEXT           NOT NULL,
    description   TEXT           NULL,
    price         NUMERIC(10, 2) NOT NULL,
    available     BOOLEAN        NOT NULL DEFAULT true,
    max_quantity  INT            NOT NULL DEFAULT 10,
    date_created  TIMESTAMP      NOT NULL,
    date_updated  TIMESTAMP      NOT NULL,

    PRIMARY KEY (addon_id),
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);

-- Version: 1.11
-- Description: Create table order_item_addons for tracking addons in orders
CREATE TABLE order_item_addons (
    order_item_addon_id UUID           NOT NULL,
    order_item_id       UUID           NOT NULL,
    addon_id            UUID           NOT NULL,
    addon_name          TEXT           NOT NULL,
    addon_price         NUMERIC(10, 2) NOT NULL,
    quantity            INT            NOT NULL,
    date_created        TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_item_addon_id),
    FOREIGN KEY (order_item_id) REFERENCES order_items(order_item_id) ON DELETE CASCADE,
    FOREIGN KEY (addon_id) REFERENCES addons(addon_id)
);

-- Version: 1.12
-- Description: Create indexes for addons
CREATE INDEX idx_addons_category_id ON addons(category_id);
CREATE INDEX idx_addons_restaurant_id ON addons(restaurant_id);
CREATE INDEX idx_order_item_addons_order_item_id ON order_item_addons(order_item_id);

-- Version: 1.13
-- Description: Add delivery location coordinates and per-restaurant delivery distance limit
ALTER TABLE restaurants
    ADD COLUMN latitude                   DOUBLE PRECISION NULL,
    ADD COLUMN longitude                  DOUBLE PRECISION NULL,
    ADD COLUMN max_delivery_distance_km   DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE delivery_addresses
    ADD COLUMN latitude  DOUBLE PRECISION NULL,
    ADD COLUMN longitude DOUBLE PRECISION NULL;

-- Version: 1.14
-- Description: Add tax_rate column to restaurants
ALTER TABLE restaurants
    ADD COLUMN tax_rate DOUBLE PRECISION NOT NULL DEFAULT 0.10;

-- Version: 1.15
-- Description: Create table promotions for promotion campaigns
CREATE TABLE promotions (
    promotion_id        UUID           NOT NULL,
    restaurant_id       UUID           NULL,
    code                TEXT           NOT NULL UNIQUE,
    name                TEXT           NOT NULL,
    description         TEXT           NULL,
    discount_type       TEXT           NOT NULL,
    discount_value      NUMERIC(10, 2) NOT NULL,
    min_order_amount    NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    max_discount_amount NUMERIC(10, 2) NULL,
    usage_limit         INT            NULL,
    usage_count         INT            NOT NULL DEFAULT 0,
    start_date          TIMESTAMP      NULL,
    end_date            TIMESTAMP      NULL,
    enabled             BOOLEAN        NOT NULL DEFAULT true,
    date_created        TIMESTAMP      NOT NULL,
    date_updated        TIMESTAMP      NOT NULL,

    PRIMARY KEY (promotion_id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);

CREATE INDEX idx_promotions_code ON promotions(code);
CREATE INDEX idx_promotions_restaurant_id ON promotions(restaurant_id);

-- Version: 1.16
-- Description: Add promo_code and discount columns to orders table
ALTER TABLE orders
    ADD COLUMN promo_code TEXT           NULL,
    ADD COLUMN discount   NUMERIC(10, 2) NOT NULL DEFAULT 0.00;

-- Version: 1.17
-- Description: Add min_spend column to restaurants
ALTER TABLE restaurants
    ADD COLUMN min_spend NUMERIC(10, 2) NOT NULL DEFAULT 0.00;

-- Version: 1.18
-- Description: Create table images for tracking uploaded image objects
CREATE TABLE images (
    image_id      UUID        NOT NULL,
    restaurant_id UUID        NOT NULL,
    entity_type   TEXT        NOT NULL,
    object_path   TEXT        NOT NULL UNIQUE,
    public_url    TEXT        NOT NULL,
    content_type  TEXT        NOT NULL,
    size_bytes    BIGINT      NOT NULL DEFAULT 0,
    status        TEXT        NOT NULL DEFAULT 'pending',
    uploaded_by   UUID        NULL,
    date_created  TIMESTAMP   NOT NULL,
    date_updated  TIMESTAMP   NOT NULL,

    PRIMARY KEY (image_id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);

CREATE INDEX idx_images_restaurant_id ON images(restaurant_id);
CREATE INDEX idx_images_status ON images(status);

-- Version: 1.19
-- Description: Create table organizations
CREATE TABLE organizations (
    organization_id UUID      NOT NULL,
    name            TEXT      NOT NULL,
    date_created    TIMESTAMP NOT NULL,
    date_updated    TIMESTAMP NOT NULL,

    PRIMARY KEY (organization_id)
);

-- Version: 1.20
-- Description: Create table organization_users
CREATE TABLE organization_users (
    organization_id UUID      NOT NULL,
    user_id         UUID      NOT NULL,
    role            TEXT      NOT NULL,
    date_created    TIMESTAMP NOT NULL,

    PRIMARY KEY (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(organization_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_organization_users_user_id ON organization_users(user_id);

-- Version: 1.21
-- Description: Add organization_id column to restaurants table
ALTER TABLE restaurants
    ADD COLUMN organization_id UUID REFERENCES organizations(organization_id) ON DELETE CASCADE;

CREATE INDEX idx_restaurants_organization_id ON restaurants(organization_id);



