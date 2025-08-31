CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cities (
  uuid           UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
  name           TEXT      NOT NULL UNIQUE
);

CREATE TABLE customers (
  uuid           UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
  name           TEXT      NOT NULL,    -- ФИО или название
  phone          TEXT      NOT NULL,
  telegram_id    BIGINT,                 -- числовой ID в Telegram
  telegram_tag   TEXT,                   -- @username
  created_at     TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE drivers (
  uuid                    UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
  name                    TEXT      NOT NULL,
  telegram_id             BIGINT   NOT NULL UNIQUE,
  telegram_tag            TEXT,                   -- @username
  notification_enabled    BOOLEAN  NOT NULL DEFAULT true,
  city_uuid               UUID      NOT NULL REFERENCES cities(uuid) ON DELETE RESTRICT,
  created_at              TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE orders (
  uuid           UUID      PRIMARY KEY DEFAULT uuid_generate_v4(),
  customer_uuid  UUID      NOT NULL REFERENCES customers(uuid) ON DELETE CASCADE,
  title          TEXT      NOT NULL,
  description    TEXT,
  weight_kg      NUMERIC   NOT NULL CHECK(weight_kg >= 0),
  length_cm      NUMERIC             CHECK(length_cm >= 0),
  width_cm       NUMERIC             CHECK(width_cm >= 0),
  height_cm      NUMERIC             CHECK(height_cm >= 0),
  from_city_uuid UUID      REFERENCES cities(uuid) ON DELETE RESTRICT,
  from_address   TEXT,
  to_city_uuid   UUID      REFERENCES cities(uuid) ON DELETE RESTRICT,
  to_address     TEXT,
  tags           TEXT[]    NOT NULL DEFAULT '{}',
  price          NUMERIC   NOT NULL CHECK(price >= 0),
  available_from DATE,
  status         TEXT      NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'archived')),
  created_at     TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_tags   ON orders USING GIN(tags);
CREATE INDEX idx_orders_price  ON orders(price);
CREATE INDEX idx_orders_weight ON orders(weight_kg);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_drivers_city  ON drivers(city_uuid);
CREATE INDEX idx_drivers_notification ON drivers(notification_enabled) WHERE notification_enabled = true; 