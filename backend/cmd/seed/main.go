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
	Slug, Name, Tagline, Description string
	Family, Gender, Concentration    string
	Image                            string
	TopNotes, HeartNotes, BaseNotes  []string
	Price                            int64
	Stock                            int
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

	_, _ = pool.Exec(context.Background(), `UPDATE products SET active=false, featured=false, updated_at=now() WHERE slug = ANY($1)`, []string{
		"nocturno-oud",
		"ambar-de-recoleta",
		"flor-de-noche",
		"citrino-frio",
		"gourmand-reserva",
		"fresco-sur",
	})

	products := []productSeed{
		{
			Slug:          "miss-armaf-chic",
			Name:          "Miss Armaf Chic",
			Tagline:       "Frutas brillantes, cítricos dulces y flores limpias",
			Description:   "Apertura de frutilla, frambuesa, pera y cítricos sobre jazmín, peonía y azahar; secado de vainilla, musk, cedro, ambroxan y musgo.",
			Family:        "Cítrico",
			Gender:        "Femenino",
			Concentration: "EDP",
			TopNotes:      []string{"Frutilla", "Frambuesa", "Pera", "Naranja", "Mandarina", "Bergamota", "Calone"},
			HeartNotes:    []string{"Jazmín", "Peonía", "Azahar"},
			BaseNotes:     []string{"Patchouli", "Musk", "Vainilla", "Ambroxan", "Cedro", "Musgo"},
			Image:         "https://armaf.com/cdn/shop/files/armaf-1design-2025-11-18T235013.146.png?v=1763491870&width=1080",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "miss-armaf-catwalk",
			Name:          "Miss Armaf Catwalk",
			Tagline:       "Mandarina fresca, durazno floral y musk dulce",
			Description:   "Un perfil cítrico floral con salida de mandarina, naranja y cítricos, corazón de durazno, jazmín y lirio, y fondo de vainilla, musk y ambroxan.",
			Family:        "Cítrico",
			Gender:        "Femenino",
			Concentration: "EDP",
			TopNotes:      []string{"Mandarina", "Naranja", "Cítricos"},
			HeartNotes:    []string{"Durazno", "Jazmín", "Lirio"},
			BaseNotes:     []string{"Vainilla", "Musk", "Ambroxan"},
			Image:         "https://armafperfume.us/cdn/shop/files/9ad438dce22642d7b50a7bd7a9c81743_tplv-omjb5zjo8w-resize-jpeg_800_800.jpg?v=1751506753&width=900",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "creme-of-clouds",
			Name:          "Creme of Clouds",
			Tagline:       "Crema batida, azúcar quemada y vainilla suave",
			Description:   "Gourmand cremoso con crema batida, azúcar tostada, leche de coco y vainilla; una estela dulce y lactónica con sensación de caramelo.",
			Family:        "Gourmand",
			Gender:        "Unisex",
			Concentration: "EDP",
			TopNotes:      []string{"Crema batida", "Azúcar quemada"},
			HeartNotes:    []string{"Leche de coco", "Vainilla"},
			BaseNotes:     []string{"Caramelo", "Musk"},
			Image:         "https://perfumeoriental.com/cdn/shop/files/creme-of-clouds-fragrance-world-edp-perfume-oriental.webp?v=1771523950&width=900",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "lattafa-eclaire",
			Name:          "Lattafa Eclaire",
			Tagline:       "Caramelo cremoso, leche y vainilla con praliné",
			Description:   "Dulce gourmand de caramelo, leche y azúcar con corazón de miel y flores blancas, apoyado en vainilla, praliné y musk.",
			Family:        "Gourmand",
			Gender:        "Unisex",
			Concentration: "EDP",
			TopNotes:      []string{"Caramelo", "Leche", "Azúcar"},
			HeartNotes:    []string{"Miel", "Flores blancas"},
			BaseNotes:     []string{"Vainilla", "Praliné", "Musk"},
			Image:         "https://www.lattafa-usa.com/cdn/shop/files/Eclaire-1_5803282e-ea5b-4de5-99a5-7d06f5cbae33.png?v=1747415649&width=1200",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "odyssey-aqua",
			Name:          "Odyssey Aqua",
			Tagline:       "Pomelo, naranja y menta sobre una base limpia amaderada",
			Description:   "Fragancia fresca de salida cítrica con pomelo, naranja y artemisia; evoluciona hacia lavanda y menta sobre ciprés, patchouli y ambroxan.",
			Family:        "Fresco",
			Gender:        "Masculino",
			Concentration: "EDP",
			TopNotes:      []string{"Pomelo", "Naranja", "Artemisia"},
			HeartNotes:    []string{"Lavanda", "Menta"},
			BaseNotes:     []string{"Ciprés", "Patchouli", "Ambroxan"},
			Image:         "https://armaf.com/cdn/shop/files/image-2023-05-04T112339.859.jpg?v=1739111570&width=1200",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "nectar-of-ecstasy-v1",
			Name:          "Nectar of Ecstasy (Versión 1)",
			Tagline:       "Fruta jugosa, bergamota y vainilla caramelada",
			Description:   "Versión dulce y cremosa de Nectar of Ecstasy, con açai, arándano y bergamota sobre flores suaves, vainilla, caramelo y musk.",
			Family:        "Gourmand",
			Gender:        "Femenino",
			Concentration: "EDP",
			TopNotes:      []string{"Açai", "Arándano", "Bergamota"},
			HeartNotes:    []string{"Fresia", "Muguet"},
			BaseNotes:     []string{"Vainilla", "Caramelo", "Musk"},
			Image:         "https://www.french-avenue-parfum.com/wp-content/uploads/2025/01/Eau-de-parfum-Nectar-of-Ecstasy.jpg",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "ameer-al-arab-imperium",
			Name:          "Ameer Al Arab Imperium",
			Tagline:       "Bergamota, jengibre y salvia con fondo amaderado",
			Description:   "Woody aromatic con apertura de bergamota, jengibre y salvia, corazón de manzana, cashmeran y geranio, y fondo de musk, ámbar y sándalo.",
			Family:        "Amaderado",
			Gender:        "Masculino",
			Concentration: "EDP",
			TopNotes:      []string{"Bergamota", "Jengibre", "Salvia"},
			HeartNotes:    []string{"Manzana", "Cashmeran", "Geranio"},
			BaseNotes:     []string{"Musk", "Ámbar", "Sándalo"},
			Image:         "https://fimgs.net/mdimg/perfume-thumbs/375x500.100148.jpg",
			Price:         2750000,
			Stock:         10,
		},
		{
			Slug:          "spicebomb-night-vision",
			Name:          "Spicebomb Night Vision",
			Tagline:       "Limón, especias negras y maderas verdes intensas",
			Description:   "EDP especiado amaderado con limón y especias negras, corazón verde resinoso con salvia e incienso, y base de abeto balsámico, cedro, patchouli y ládano.",
			Family:        "Amaderado",
			Gender:        "Masculino",
			Concentration: "EDP",
			TopNotes:      []string{"Limón", "Pimienta negra", "Chile negro", "Nuez moscada", "Clavo"},
			HeartNotes:    []string{"Salvia esclarea", "Resina verde", "Incienso"},
			BaseNotes:     []string{"Abeto balsámico", "Cedro", "Patchouli", "Ládano"},
			Image:         "https://us.viktor-rolf.com/dw/image/v2/AANG_PRD/on/demandware.static/-/Sites-vr-master-catalog/default/dw06340d3b/SB%20NV%20EDP%202024/01_vr_frag_spb_night_vision_edp_perfect_pdp_premium_packshot_90ml_1x1.jpg?q=70&sfrm=jpg&sh=900&sm=cut&sw=900",
			Price:         2750000,
			Stock:         10,
		},
		{
			Slug:          "nectar-of-ecstasy-v2",
			Name:          "Nectar of Ecstasy (Versión 2)",
			Tagline:       "Cítricos dulces con firma de cedro y ámbar",
			Description:   "Versión amaderada y dulce de Nectar of Ecstasy, con fruta roja, bergamota y flores claras sobre cedro, ámbar y notas dulces.",
			Family:        "Amaderado",
			Gender:        "Femenino",
			Concentration: "EDP",
			TopNotes:      []string{"Açai", "Arándano", "Bergamota"},
			HeartNotes:    []string{"Fresia", "Muguet"},
			BaseNotes:     []string{"Cedro", "Ámbar", "Notas dulces"},
			Image:         "https://www.french-avenue-parfum.com/wp-content/uploads/2025/01/Eau-de-parfum-Nectar-of-Ecstasy.jpg",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "odyssey-limoni",
			Name:          "Odyssey Limoni",
			Tagline:       "Limón dulce, naranja y té con frescura marina",
			Description:   "Cítrico fresco con limón, naranja dulce, mandarina y bergamota; corazón de flor de azahar, notas marinas y jengibre, y base de té, musk y ámbar.",
			Family:        "Cítrico",
			Gender:        "Unisex",
			Concentration: "EDP",
			TopNotes:      []string{"Limón", "Naranja dulce", "Mandarina", "Bergamota"},
			HeartNotes:    []string{"Azahar", "Notas marinas", "Jengibre"},
			BaseNotes:     []string{"Té", "Musk", "Ámbar"},
			Image:         "https://fimgs.net/mdimg/perfume-thumbs/375x500.98695.jpg",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "odyssey-aqua-v2",
			Name:          "Odyssey Aqua (Versión 2)",
			Tagline:       "Menta fresca, lavanda y cítricos acuáticos",
			Description:   "Segunda lectura de Odyssey Aqua centrada en frescura y menta, con pomelo, naranja, artemisia, lavanda y un cierre de ambroxan, patchouli y ciprés.",
			Family:        "Fresco",
			Gender:        "Masculino",
			Concentration: "EDP",
			TopNotes:      []string{"Pomelo", "Naranja", "Artemisia"},
			HeartNotes:    []string{"Menta", "Lavanda"},
			BaseNotes:     []string{"Ambroxan", "Patchouli", "Ciprés"},
			Image:         "https://armaf.com/cdn/shop/files/IMG_20250710_1558211-ezgif.com-webp-to-jpg-converter.jpg?v=1767894320&width=1200",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "paris-corner-taskeen-dia",
			Name:          "Paris Corner Taskeen Día",
			Tagline:       "Durazno, naranja sanguina y vainilla dulce",
			Description:   "Perfil frutal dulce de durazno, naranja sanguina y cardamomo, con corazón de heliotropo, davana, cognac y jazmín, y base de sándalo, vainilla, tonka y patchouli.",
			Family:        "Gourmand",
			Gender:        "Unisex",
			Concentration: "EDP",
			TopNotes:      []string{"Durazno", "Naranja sanguina", "Cardamomo"},
			HeartNotes:    []string{"Heliotropo", "Davana", "Cognac", "Jazmín"},
			BaseNotes:     []string{"Sándalo", "Benjuí", "Cashmeran", "Vainilla", "Tonka", "Ládano", "Patchouli"},
			Image:         "https://www.pariscornerperfumes.com/cdn/shop/products/s-l1600_8555abe4-8c06-4e28-9072-d21df9f08a1e.jpg?v=1644237892&width=900",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "paris-corner-taskeen-noche",
			Name:          "Paris Corner Taskeen Noche",
			Tagline:       "Fruta dulce, cognac suave y fondo cremoso",
			Description:   "Perfil nocturno de Taskeen con durazno y naranja sanguina, un corazón floral con davana y cognac, y una base dulce de sándalo, benjuí, vainilla, tonka, ládano y patchouli.",
			Family:        "Gourmand",
			Gender:        "Unisex",
			Concentration: "EDP",
			TopNotes:      []string{"Durazno", "Naranja sanguina", "Cardamomo"},
			HeartNotes:    []string{"Heliotropo", "Davana", "Cognac", "Jazmín"},
			BaseNotes:     []string{"Sándalo", "Benjuí", "Cashmeran", "Vainilla", "Tonka", "Ládano", "Patchouli"},
			Image:         "https://www.pariscornerperfumes.com/cdn/shop/products/s-l1600_8555abe4-8c06-4e28-9072-d21df9f08a1e.jpg?v=1644237892&width=900",
			Price:         2950000,
			Stock:         10,
		},
		{
			Slug:          "paris-corner-khair-confection",
			Name:          "Paris Corner Khair Confection",
			Tagline:       "Pera, crema batida y vainilla malvavisco",
			Description:   "Gourmand dulce y frutal con pera y crema batida, corazón de jazmín, ylang-ylang y cashmeran, y fondo de sándalo, malvavisco y vainilla.",
			Family:        "Gourmand",
			Gender:        "Unisex",
			Concentration: "EDP",
			TopNotes:      []string{"Pera", "Crema batida"},
			HeartNotes:    []string{"Jazmín", "Ylang-ylang", "Cashmeran"},
			BaseNotes:     []string{"Sándalo", "Malvavisco", "Vainilla"},
			Image:         "https://www.pariscornerperfumes.com/cdn/shop/files/KHAIRCONFECTION01.jpg?v=1726730340&width=900",
			Price:         2950000,
			Stock:         10,
		},
	}
	for i, p := range products {
		var id string
		err := pool.QueryRow(context.Background(), `INSERT INTO products (slug,name,tagline,description,scent_family,gender_tag,concentration,top_notes,heart_notes,base_notes,featured,active,display_order)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,true,true,$11)
ON CONFLICT (slug) DO UPDATE SET
	name=EXCLUDED.name,
	tagline=EXCLUDED.tagline,
	description=EXCLUDED.description,
	scent_family=EXCLUDED.scent_family,
	gender_tag=EXCLUDED.gender_tag,
	concentration=EXCLUDED.concentration,
	top_notes=EXCLUDED.top_notes,
	heart_notes=EXCLUDED.heart_notes,
	base_notes=EXCLUDED.base_notes,
	featured=EXCLUDED.featured,
	active=EXCLUDED.active,
	display_order=EXCLUDED.display_order,
	updated_at=now()
RETURNING id`,
			p.Slug, p.Name, p.Tagline, p.Description, p.Family, p.Gender, p.Concentration, p.TopNotes, p.HeartNotes, p.BaseNotes, i).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		sku := p.Slug + "-100"
		_, _ = pool.Exec(context.Background(), `INSERT INTO product_variants (product_id,size_ml,price_ars_cents,stock,sku,active,weight_grams)
VALUES ($1,100,$2,$3,$4,true,320)
ON CONFLICT (sku) DO UPDATE SET price_ars_cents=EXCLUDED.price_ars_cents, stock=EXCLUDED.stock, active=true, weight_grams=EXCLUDED.weight_grams
WHERE product_variants.product_id=EXCLUDED.product_id`, id, p.Price, p.Stock, sku)
		_, _ = pool.Exec(context.Background(), `UPDATE product_variants SET active=false WHERE product_id=$1 AND sku<>$2`, id, sku)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_images WHERE product_id=$1`, id)
		_, _ = pool.Exec(context.Background(), `INSERT INTO product_images (product_id,url,alt_text,is_primary,sort_order) VALUES ($1,$2,$3,true,0)`, id, p.Image, p.Name)
	}
	log.Println("seed data applied; change admin password immediately")
}
