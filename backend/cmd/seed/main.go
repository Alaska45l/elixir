package main

import (
	"context"
	"log"
	"strings"

	"elixir/backend/internal/config"
	"elixir/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type productSeed struct {
	Slug, Name, Tagline, Family, Image string
	Price                              int64
	Stock                              int
}

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" || strings.Contains(cfg.DatabaseURL, "user:pass@host/dbname") {
		log.Fatal("DATABASE_URL is not configured; edit backend/.env with your Neon/PostgreSQL connection string")
	}
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if pool != nil {
		defer pool.Close()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("elixir2024"), 12)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = pool.Exec(context.Background(), `INSERT INTO admin_users (username,password_hash) VALUES ('admin',$1) ON CONFLICT (username) DO UPDATE SET password_hash=EXCLUDED.password_hash`, string(hash))
	_, _ = pool.Exec(context.Background(), `INSERT INTO homepage_settings (id, hero_heading, hero_subheading, hero_image_url, hero_cta_label, hero_cta_url, editorial_heading, editorial_body, editorial_image_url) VALUES (1,'Perfumería argentina de gesto privado','Fragancias intensas, precisas y comerciales, curadas para noches largas, hoteles silenciosos y piel con presencia.','https://images.unsplash.com/photo-1619994403073-2cec844b8e63?auto=format&fit=crop&w=1200&q=85','Catálogo','/fragrances','Una firma de baja voz','ELIXIR Exclusive trabaja maderas limpias, ámbar seco, flores oscuras y cítricos fríos. Cada compra se prepara con empaque sobrio y seguimiento personalizado.','https://images.unsplash.com/photo-1595425970377-c9703cf48b6f?auto=format&fit=crop&w=1000&q=85') ON CONFLICT (id) DO NOTHING`)
	_, _ = pool.Exec(context.Background(), `INSERT INTO shipping_zones (zone_name,province_codes,base_cost_cents,estimated_days_min,estimated_days_max) VALUES ('CABA',ARRAY['CF'],0,0,3),('Gran Buenos Aires',ARRAY['BA'],250000,1,4),('Interior',ARRAY['AR'],420000,3,7) ON CONFLICT DO NOTHING`)
	_, _ = pool.Exec(context.Background(), `INSERT INTO discount_codes (code,discount_type,discount_value,min_order_cents,active) VALUES ('BIENVENIDO10','percent',10,0,true),('ELIXIR200','fixed',20000,0,true) ON CONFLICT (code) DO NOTHING`)

	products := []productSeed{
		{"nocturno-oud", "Nocturno Oud", "Oud seco, rosa negra y cuero limpio", "Amaderado", "https://images.unsplash.com/photo-1592945403244-b3fbafd7f539?auto=format&fit=crop&w=900&q=85", 8900000, 4},
		{"ambar-de-recoleta", "Ámbar de Recoleta", "Ámbar cálido, vainilla sobria e incienso", "Oriental", "https://images.unsplash.com/photo-1615634260167-c8cdede054de?auto=format&fit=crop&w=900&q=85", 7600000, 12},
		{"flor-de-noche", "Flor de Noche", "Jazmín oscuro, iris y almizcle limpio", "Floral", "https://images.unsplash.com/photo-1547887538-e3a2f32cb1cc?auto=format&fit=crop&w=900&q=85", 8200000, 3},
		{"citrino-frio", "Citrino Frío", "Bergamota helada, neroli y cedro blanco", "Cítrico", "https://images.unsplash.com/photo-1600612253971-422e7f7faeb6?auto=format&fit=crop&w=900&q=85", 6900000, 8},
		{"gourmand-reserva", "Gourmand Reserva", "Tonka, cacao amargo y sándalo", "Gourmand", "https://images.unsplash.com/photo-1594035910387-fea47794261f?auto=format&fit=crop&w=900&q=85", 9300000, 2},
		{"fresco-sur", "Fresco Sur", "Mate verde, pomelo y vetiver", "Fresco", "https://images.unsplash.com/photo-1585386959984-a4155223168f?auto=format&fit=crop&w=900&q=85", 7100000, 0},
	}
	for i, p := range products {
		var id string
		err := pool.QueryRow(context.Background(), `INSERT INTO products (slug,name,tagline,description,scent_family,gender_tag,concentration,top_notes,heart_notes,base_notes,featured,display_order) VALUES ($1,$2,$3,$4,$5,'Unisex','EDP',ARRAY['Bergamota','Pimienta rosa','Azafrán'],ARRAY['Rosa','Iris','Incienso'],ARRAY['Sándalo','Ámbar','Almizcle'],true,$6) ON CONFLICT (slug) DO UPDATE SET name=EXCLUDED.name RETURNING id`,
			p.Slug, p.Name, p.Tagline, p.Tagline+". Una composición de alta permanencia pensada para piel y clima urbano argentino.", p.Family, i).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = pool.Exec(context.Background(), `INSERT INTO product_variants (product_id,size_ml,price_ars_cents,stock,sku) VALUES ($1,50,$2,$3,$4),($1,100,$5,$3,$6) ON CONFLICT (sku) DO NOTHING`, id, p.Price, p.Stock, p.Slug+"-50", p.Price*165/100, p.Slug+"-100")
		_, _ = pool.Exec(context.Background(), `INSERT INTO product_images (product_id,url,alt_text,is_primary,sort_order) VALUES ($1,$2,$3,true,0),($1,$4,$5,false,1)`, id, p.Image, p.Name, p.Image, p.Name+" detalle")
	}
	log.Println("seed data applied; change admin password immediately")
}
