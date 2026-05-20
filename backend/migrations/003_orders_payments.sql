CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_reference TEXT UNIQUE NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  customer_name TEXT NOT NULL,
  customer_email TEXT NOT NULL,
  customer_phone TEXT,
  shipping_address JSONB,
  shipping_cost_ars_cents BIGINT DEFAULT 0,
  subtotal_ars_cents BIGINT NOT NULL,
  total_ars_cents BIGINT NOT NULL,
  discount_code TEXT,
  discount_ars_cents BIGINT DEFAULT 0,
  currency TEXT DEFAULT 'ARS',
  tracking_number TEXT,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT orders_status_check CHECK (status IN ('pending', 'paid', 'failed', 'cancelled', 'shipped', 'delivered'))
);

CREATE TABLE IF NOT EXISTS order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
  product_id UUID REFERENCES products(id),
  variant_id UUID REFERENCES product_variants(id),
  product_name TEXT NOT NULL,
  size_ml INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  unit_price_ars_cents BIGINT NOT NULL,
  subtotal_ars_cents BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS payment_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID REFERENCES orders(id),
  mp_payment_id TEXT UNIQUE NOT NULL,
  mp_preference_id TEXT,
  mp_status TEXT,
  mp_status_detail TEXT,
  raw_payload JSONB,
  processed_at TIMESTAMPTZ DEFAULT now()
);
