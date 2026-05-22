ALTER TABLE product_variants
  ADD COLUMN IF NOT EXISTS weight_grams INTEGER NOT NULL DEFAULT 200;

CREATE TABLE IF NOT EXISTS site_settings (
  id INTEGER PRIMARY KEY DEFAULT 1,
  footer_instagram_url TEXT DEFAULT '',
  footer_tiktok_url TEXT DEFAULT '',
  footer_whatsapp_url TEXT DEFAULT '',
  announcement_bar_text TEXT DEFAULT 'Envíos a todo el país · Empaque discreto · Seguimiento personalizado',
  announcement_bar_active BOOLEAN DEFAULT true,
  about_title TEXT DEFAULT 'ELIXIR Exclusive',
  about_description TEXT DEFAULT 'Perfumería argentina de lujo discreto. Fragancias intensas, envíos nacionales y atención privada.',
  about_location TEXT DEFAULT 'Buenos Aires, Argentina',
  about_phone TEXT DEFAULT '',
  faq_items JSONB DEFAULT '[]'::jsonb,
  return_policy_html TEXT DEFAULT '',
  navbar_product_categories JSONB DEFAULT '[]'::jsonb,
  updated_at TIMESTAMPTZ DEFAULT now(),
  CONSTRAINT site_settings_singleton CHECK (id = 1)
);

INSERT INTO site_settings (
  id,
  footer_instagram_url,
  footer_tiktok_url,
  footer_whatsapp_url,
  announcement_bar_text,
  announcement_bar_active,
  about_title,
  about_description,
  about_location,
  about_phone,
  faq_items,
  return_policy_html,
  navbar_product_categories
) VALUES (
  1,
  '',
  '',
  '',
  'Envíos a todo el país · Empaque discreto · Seguimiento personalizado',
  true,
  'ELIXIR Exclusive',
  'Perfumería argentina de lujo discreto. Fragancias intensas, envíos nacionales y atención privada.',
  'Buenos Aires, Argentina',
  '',
  '[
    {"question":"¿Los perfumes son originales?","answer":"Sí. ELIXIR Exclusive comercializa fragancias seleccionadas y documentadas."},
    {"question":"¿Qué medios de pago aceptan?","answer":"El checkout opera en ARS mediante MercadoPago."},
    {"question":"¿Hacen envíos?","answer":"Sí, a CABA, GBA e Interior con seguimiento."},
    {"question":"¿Puedo consultar por WhatsApp?","answer":"Sí. Recomendamos WhatsApp para asesoramiento rápido."}
  ]'::jsonb,
  '<p>Los cambios se revisan caso por caso con el producto cerrado, sin uso y dentro de los plazos informados por atención al cliente.</p>',
  '[
    {"label":"Fragancias Masculinas","href":"/fragrances?gender=Masculino"},
    {"label":"Fragancias Femeninas","href":"/fragrances?gender=Femenino"},
    {"label":"Línea Oriental","href":"/fragrances?family=Oriental"},
    {"label":"Línea Amaderada","href":"/fragrances?family=Amaderado"}
  ]'::jsonb
) ON CONFLICT (id) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
