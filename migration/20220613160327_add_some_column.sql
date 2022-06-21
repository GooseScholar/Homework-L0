-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE orders
(
    order_uid VARCHAR(32),
    track_number VARCHAR(32),
    "entry" VARCHAR(16),
    locale VARCHAR(16),
    internal_signature VARCHAR(64),
    customer_id VARCHAR(32),
    delivery_service VARCHAR(32),
    shardkey VARCHAR(16),
    sm_id int NOT NULL,
    date_created VARCHAR(32),
    oof_shard VARCHAR(32)
);

CREATE TABLE delivery
(
    order_uid VARCHAR(32),
    "name" VARCHAR(32),
    phone VARCHAR(16),
    zip VARCHAR(16),
    city VARCHAR(32),
    "address" VARCHAR(32),
    region VARCHAR(16),
    email VARCHAR(32)
);

CREATE TABLE payment
(
    order_uid VARCHAR(32),
    "transaction" VARCHAR(32),
    request_id VARCHAR(32),
    currency VARCHAR(16),
    "provider" VARCHAR(16),
    amount int NOT NULL,
    payment_dt bigint NOT NULL,
    bank VARCHAR(16),
    delivery_cost int NOT NULL,
    goods_total int NOT NULL,
    custom_fee int NOT NULL
);

CREATE TABLE items
(
    order_uid VARCHAR(32),
    chrt_id int NOT NULL,
    track_number VARCHAR(32),
    price int NOT NULL,
    rid VARCHAR(32),
    "name" VARCHAR(32),
    sale int NOT NULL,
    size VARCHAR(16),
    total_price int NOT NULL,
    nm_id int NOT NULL,
    brand VARCHAR(32),
    "status" int NOT NULL
);




-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE orders;

CREATE TABLE delivery;

DROP TABLE payment;

DROP TABLE items;