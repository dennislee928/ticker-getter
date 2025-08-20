-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- 用戶索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_tls_fingerprint ON users(tls_fingerprint);

-- 活動索引
CREATE INDEX IF NOT EXISTS idx_events_created_by ON events(created_by);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
CREATE INDEX IF NOT EXISTS idx_events_end_time ON events(end_time);

-- 票種索引
CREATE INDEX IF NOT EXISTS idx_ticket_types_event_id ON ticket_types(event_id);
CREATE INDEX IF NOT EXISTS idx_ticket_types_sale_start ON ticket_types(sale_start);
CREATE INDEX IF NOT EXISTS idx_ticket_types_sale_end ON ticket_types(sale_end);

-- 訂單索引
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status);

-- 訂單項目索引
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_ticket_type_id ON order_items(ticket_type_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_tls_fingerprint;

DROP INDEX IF EXISTS idx_events_created_by;
DROP INDEX IF EXISTS idx_events_start_time;
DROP INDEX IF EXISTS idx_events_end_time;

DROP INDEX IF EXISTS idx_ticket_types_event_id;
DROP INDEX IF EXISTS idx_ticket_types_sale_start;
DROP INDEX IF EXISTS idx_ticket_types_sale_end;

DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_payment_status;

DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_order_items_ticket_type_id;
