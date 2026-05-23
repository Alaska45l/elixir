ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS shipping_carrier TEXT,
  ADD COLUMN IF NOT EXISTS shipped_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS internal_notes TEXT;

ALTER TABLE site_settings
  ADD COLUMN IF NOT EXISTS low_stock_threshold INTEGER NOT NULL DEFAULT 5;

CREATE INDEX IF NOT EXISTS idx_payment_events_order_processed ON payment_events(order_id, processed_at DESC);
