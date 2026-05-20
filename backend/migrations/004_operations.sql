CREATE TABLE IF NOT EXISTS discount_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  discount_type TEXT NOT NULL,
  discount_value BIGINT NOT NULL,
  min_order_cents BIGINT DEFAULT 0,
  max_uses INTEGER,
  uses INTEGER DEFAULT 0,
  active BOOLEAN DEFAULT true,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT discount_type_check CHECK (discount_type IN ('percent', 'fixed'))
);

CREATE TABLE IF NOT EXISTS contact_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  subject TEXT,
  message TEXT NOT NULL,
  read BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS shipping_zones (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  zone_name TEXT NOT NULL,
  province_codes TEXT[],
  base_cost_cents BIGINT NOT NULL,
  per_kg_cents BIGINT DEFAULT 0,
  estimated_days_min INTEGER,
  estimated_days_max INTEGER,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS abandoned_carts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT,
  cart_data JSONB,
  created_at TIMESTAMPTZ DEFAULT now(),
  recovered BOOLEAN DEFAULT false
);
