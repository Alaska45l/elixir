ALTER TABLE homepage_settings
  ADD COLUMN IF NOT EXISTS hero_image_mode TEXT DEFAULT 'product_covers',
  ADD COLUMN IF NOT EXISTS hero_rotation_interval_ms INTEGER DEFAULT 8000;

UPDATE homepage_settings
SET
  hero_image_mode = CASE
    WHEN hero_image_mode IN ('static', 'product_covers') THEN hero_image_mode
    ELSE 'product_covers'
  END,
  hero_rotation_interval_ms = CASE
    WHEN hero_rotation_interval_ms BETWEEN 1000 AND 60000 THEN hero_rotation_interval_ms
    ELSE 8000
  END;

ALTER TABLE homepage_settings
  ALTER COLUMN hero_image_mode SET NOT NULL,
  ALTER COLUMN hero_rotation_interval_ms SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'homepage_hero_image_mode_check'
  ) THEN
    ALTER TABLE homepage_settings
      ADD CONSTRAINT homepage_hero_image_mode_check
      CHECK (hero_image_mode IN ('static', 'product_covers'));
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'homepage_hero_rotation_interval_check'
  ) THEN
    ALTER TABLE homepage_settings
      ADD CONSTRAINT homepage_hero_rotation_interval_check
      CHECK (hero_rotation_interval_ms BETWEEN 1000 AND 60000);
  END IF;
END $$;
