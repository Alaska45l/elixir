CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS admin_users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  last_login_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS homepage_settings (
  id INTEGER PRIMARY KEY DEFAULT 1,
  hero_heading TEXT,
  hero_subheading TEXT,
  hero_image_url TEXT,
  hero_cta_label TEXT,
  hero_cta_url TEXT,
  editorial_heading TEXT,
  editorial_body TEXT,
  editorial_image_url TEXT,
  updated_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT homepage_singleton CHECK (id = 1)
);
