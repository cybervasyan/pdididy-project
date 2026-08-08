-- +goose Up
CREATE table orders
(
    order_uuid       UUID PRIMARY KEY,
    user_uuid        UUID NOT NULL,
	part_uuids       UUID[] NOT NULL,
	total_price      DOUBLE PRECISION NOT NULL,
	transaction_uuid UUID,
	payment_method   text,
	status          text NOT NULL
);

-- +goose Down
drop table if exists orders;
