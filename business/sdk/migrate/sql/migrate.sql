-- Version: 1.01
-- Description: Create Food Flow schema

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

CREATE TABLE organizations (
    organization_id UUID      NOT NULL,
    name            TEXT      NOT NULL,
    date_created    TIMESTAMP NOT NULL,
    date_updated    TIMESTAMP NOT NULL,

    PRIMARY KEY (organization_id)
);

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

CREATE TABLE restaurants (
    restaurant_id            UUID             NOT NULL,
    organization_id          UUID             NULL,
    name                     TEXT             NOT NULL,
    description              TEXT             NULL,
    address                  TEXT             NOT NULL,
    phone                    TEXT             NOT NULL,
    email                    TEXT             NOT NULL,
    image_url                TEXT             NULL,
    logo_url                 TEXT             NOT NULL DEFAULT '',
    operating_hours          JSONB            NOT NULL DEFAULT '{
        "monday":    {"open": "10:00", "close": "22:00", "isClosed": false},
        "tuesday":   {"open": "10:00", "close": "22:00", "isClosed": false},
        "wednesday": {"open": "10:00", "close": "22:00", "isClosed": false},
        "thursday":  {"open": "10:00", "close": "22:00", "isClosed": false},
        "friday":    {"open": "10:00", "close": "23:00", "isClosed": false},
        "saturday":  {"open": "11:00", "close": "23:00", "isClosed": false},
        "sunday":    {"open": "11:00", "close": "22:00", "isClosed": false}
    }'::jsonb,
    enabled                  BOOLEAN          NOT NULL,
    latitude                 DOUBLE PRECISION NULL,
    longitude                DOUBLE PRECISION NULL,
    max_delivery_distance_km DOUBLE PRECISION NOT NULL DEFAULT 0,
    tax_rate                 DOUBLE PRECISION NOT NULL DEFAULT 0.10,
    min_spend                NUMERIC(10, 2)   NOT NULL DEFAULT 0.00,
    date_created             TIMESTAMP        NOT NULL,
    date_updated             TIMESTAMP        NOT NULL,

    PRIMARY KEY (restaurant_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(organization_id) ON DELETE CASCADE
);

CREATE INDEX idx_restaurants_organization_id ON restaurants(organization_id);

CREATE TABLE categories (
    category_id   UUID      NOT NULL,
    name          TEXT      NOT NULL,
    description   TEXT      NULL,
    restaurant_id UUID      NOT NULL,
    enabled       BOOLEAN   NOT NULL,
    rank          INT       NULL,
    date_created  TIMESTAMP NOT NULL,
    date_updated  TIMESTAMP NOT NULL,

    PRIMARY KEY (category_id),
    CONSTRAINT categories_id_restaurant_unique
        UNIQUE (category_id, restaurant_id),
    FOREIGN KEY (restaurant_id)
        REFERENCES restaurants(restaurant_id)
        ON DELETE CASCADE,
    CONSTRAINT categories_rank_check
        CHECK (rank IS NULL OR rank >= 1)
);

CREATE INDEX idx_categories_restaurant_rank
    ON categories(restaurant_id, rank, name, category_id);

CREATE TABLE menu_items (
    menu_item_id  UUID           NOT NULL,
    name          TEXT           NOT NULL,
    description   TEXT           NULL,
    price         NUMERIC(10, 2) NOT NULL,
    category_id   UUID           NOT NULL,
    restaurant_id UUID           NOT NULL,
    image_url     TEXT           NULL,
    available     BOOLEAN        NOT NULL,
    rank          INT            NULL,
    date_created  TIMESTAMP      NOT NULL,
    date_updated  TIMESTAMP      NOT NULL,

    PRIMARY KEY (menu_item_id),
    CONSTRAINT menu_items_id_restaurant_unique
        UNIQUE (menu_item_id, restaurant_id),
    CONSTRAINT menu_items_category_restaurant_fkey
        FOREIGN KEY (category_id, restaurant_id)
        REFERENCES categories(category_id, restaurant_id)
        ON DELETE CASCADE,
    FOREIGN KEY (restaurant_id)
        REFERENCES restaurants(restaurant_id)
        ON DELETE CASCADE,
    CONSTRAINT menu_items_price_check
        CHECK (price > 0),
    CONSTRAINT menu_items_rank_check
        CHECK (rank IS NULL OR rank >= 1)
);

CREATE INDEX idx_menu_items_restaurant_rank
    ON menu_items(restaurant_id, rank, name, menu_item_id);
CREATE INDEX idx_menu_items_category_rank
    ON menu_items(category_id, rank, name, menu_item_id);

CREATE TABLE modifier_groups (
    modifier_group_id UUID      NOT NULL,
    menu_item_id      UUID      NOT NULL,
    restaurant_id     UUID      NOT NULL,
    name              TEXT      NOT NULL,
    description       TEXT      NULL,
    min_selections    INT       NOT NULL DEFAULT 1,
    max_selections    INT       NOT NULL DEFAULT 1,
    available         BOOLEAN   NOT NULL DEFAULT false,
    rank              INT       NULL,
    date_created      TIMESTAMP NOT NULL,
    date_updated      TIMESTAMP NOT NULL,

    PRIMARY KEY (modifier_group_id),
    UNIQUE (modifier_group_id, restaurant_id),
    FOREIGN KEY (menu_item_id, restaurant_id)
        REFERENCES menu_items(menu_item_id, restaurant_id)
        ON DELETE CASCADE,
    CHECK (min_selections IN (0, 1)),
    CHECK (max_selections = 1),
    CHECK (min_selections <= max_selections),
    CHECK (rank IS NULL OR rank >= 1)
);

CREATE UNIQUE INDEX idx_modifier_groups_unique_name
    ON modifier_groups(menu_item_id, lower(name));
CREATE INDEX idx_modifier_groups_menu_item_rank
    ON modifier_groups(menu_item_id, rank, name, modifier_group_id);
CREATE INDEX idx_modifier_groups_restaurant
    ON modifier_groups(restaurant_id);

CREATE TABLE modifier_options (
    modifier_option_id UUID           NOT NULL,
    modifier_group_id  UUID           NOT NULL,
    restaurant_id      UUID           NOT NULL,
    name               TEXT           NOT NULL,
    description        TEXT           NULL,
    price_delta        NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    available          BOOLEAN        NOT NULL DEFAULT true,
    rank               INT            NULL,
    date_created       TIMESTAMP      NOT NULL,
    date_updated       TIMESTAMP      NOT NULL,

    PRIMARY KEY (modifier_option_id),
    UNIQUE (modifier_option_id, restaurant_id),
    FOREIGN KEY (modifier_group_id, restaurant_id)
        REFERENCES modifier_groups(modifier_group_id, restaurant_id)
        ON DELETE CASCADE,
    CHECK (price_delta >= 0),
    CHECK (rank IS NULL OR rank >= 1)
);

CREATE UNIQUE INDEX idx_modifier_options_unique_name
    ON modifier_options(modifier_group_id, lower(name));
CREATE INDEX idx_modifier_options_group_rank
    ON modifier_options(modifier_group_id, rank, name, modifier_option_id);
CREATE INDEX idx_modifier_options_restaurant
    ON modifier_options(restaurant_id);

CREATE TABLE addons (
    addon_id      UUID           NOT NULL,
    menu_item_id  UUID           NOT NULL,
    restaurant_id UUID           NOT NULL,
    name          TEXT           NOT NULL,
    description   TEXT           NULL,
    price         NUMERIC(10, 2) NOT NULL,
    available     BOOLEAN        NOT NULL DEFAULT true,
    max_quantity  INT            NOT NULL DEFAULT 10,
    rank          INT            NULL,
    date_created  TIMESTAMP      NOT NULL,
    date_updated  TIMESTAMP      NOT NULL,

    PRIMARY KEY (addon_id),
    CONSTRAINT addons_id_restaurant_unique
        UNIQUE (addon_id, restaurant_id),
    FOREIGN KEY (menu_item_id, restaurant_id)
        REFERENCES menu_items(menu_item_id, restaurant_id)
        ON DELETE CASCADE,
    FOREIGN KEY (restaurant_id)
        REFERENCES restaurants(restaurant_id)
        ON DELETE CASCADE,
    CONSTRAINT addons_price_check
        CHECK (price >= 0),
    CONSTRAINT addons_max_quantity_check
        CHECK (max_quantity >= 1),
    CONSTRAINT addons_rank_check
        CHECK (rank IS NULL OR rank >= 1)
);

CREATE UNIQUE INDEX idx_addons_unique_name
    ON addons(menu_item_id, lower(name));
CREATE INDEX idx_addons_menu_item_rank
    ON addons(menu_item_id, rank, name, addon_id);
CREATE INDEX idx_addons_restaurant
    ON addons(restaurant_id);

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
    discount                 NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    delivery_fee             NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    tax                      NUMERIC(10, 2) NOT NULL,
    total                    NUMERIC(10, 2) NOT NULL,
    special_instructions     TEXT           NULL,
    promo_code               TEXT           NULL,
    stripe_payment_intent_id TEXT           NULL,
    date_created             TIMESTAMP      NOT NULL,
    date_updated             TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);

CREATE INDEX idx_orders_restaurant_id ON orders(restaurant_id);
CREATE INDEX idx_orders_customer_email ON orders(customer_email);
CREATE INDEX idx_orders_order_status ON orders(order_status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_date_created ON orders(date_created DESC);

CREATE TABLE order_items (
    order_item_id        UUID           NOT NULL,
    order_id             UUID           NOT NULL,
    category_id          UUID           NOT NULL,
    category_name        TEXT           NOT NULL,
    menu_item_id         UUID           NOT NULL,
    menu_item_name       TEXT           NOT NULL,
    menu_item_price      NUMERIC(10, 2) NOT NULL,
    quantity             INT            NOT NULL,
    special_instructions TEXT           NULL,
    date_created         TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_item_id),
    FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);

CREATE TABLE order_item_addons (
    order_item_addon_id UUID           NOT NULL,
    order_item_id       UUID           NOT NULL,
    addon_id            UUID           NOT NULL,
    addon_name          TEXT           NOT NULL,
    addon_price         NUMERIC(10, 2) NOT NULL,
    quantity            INT            NOT NULL,
    date_created        TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_item_addon_id),
    FOREIGN KEY (order_item_id)
        REFERENCES order_items(order_item_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_order_item_addons_order_item_id ON order_item_addons(order_item_id);

CREATE TABLE order_item_modifiers (
    order_item_modifier_id UUID           NOT NULL,
    order_item_id          UUID           NOT NULL,
    modifier_group_id      UUID           NOT NULL,
    modifier_group_name    TEXT           NOT NULL,
    modifier_option_id     UUID           NOT NULL,
    modifier_option_name   TEXT           NOT NULL,
    price_delta            NUMERIC(10, 2) NOT NULL,
    date_created           TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_item_modifier_id),
    UNIQUE (order_item_id, modifier_group_id),
    FOREIGN KEY (order_item_id)
        REFERENCES order_items(order_item_id)
        ON DELETE CASCADE,
    CHECK (price_delta >= 0)
);

CREATE INDEX idx_order_item_modifiers_order_item
    ON order_item_modifiers(order_item_id);
CREATE INDEX idx_order_item_modifiers_option
    ON order_item_modifiers(modifier_option_id);

CREATE TABLE delivery_addresses (
    address_id            UUID             NOT NULL,
    order_id              UUID             NOT NULL UNIQUE,
    street                TEXT             NOT NULL,
    city                  TEXT             NOT NULL,
    state                 TEXT             NOT NULL,
    postal_code           TEXT             NOT NULL,
    delivery_instructions TEXT             NULL,
    latitude              DOUBLE PRECISION NULL,
    longitude             DOUBLE PRECISION NULL,
    date_created          TIMESTAMP        NOT NULL,

    PRIMARY KEY (address_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

CREATE INDEX idx_delivery_addresses_order_id ON delivery_addresses(order_id);

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

CREATE INDEX idx_payment_transactions_order_id ON payment_transactions(order_id);
CREATE INDEX idx_payment_transactions_stripe_id ON payment_transactions(stripe_payment_intent_id);

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
