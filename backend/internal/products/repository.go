package products

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func (r Repository) List(ctx context.Context, f ListFilters) (ListResult, error) {
	if r.Pool == nil {
		return ListResult{Items: []Product{}, Limit: f.Limit, Offset: f.Offset}, nil
	}
	limit := f.Limit
	if limit <= 0 || limit > 60 {
		limit = 24
	}
	args := []any{}
	where := []string{"p.active = true"}
	next := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}
	if f.Featured != nil {
		where = append(where, "p.featured = "+next(*f.Featured))
	}
	if f.Family != "" {
		where = append(where, "p.scent_family = "+next(f.Family))
	}
	if f.Gender != "" {
		where = append(where, "p.gender_tag = "+next(f.Gender))
	}
	if f.Concentration != "" {
		where = append(where, "p.concentration = "+next(f.Concentration))
	}
	if f.Search != "" {
		where = append(where, "(p.name ILIKE "+next("%"+f.Search+"%")+" OR p.tagline ILIKE "+next("%"+f.Search+"%")+")")
	}
	if f.InStock {
		where = append(where, "EXISTS (SELECT 1 FROM product_variants v WHERE v.product_id=p.id AND v.active=true AND v.stock > 0)")
	}
	if f.MinPrice > 0 {
		where = append(where, "EXISTS (SELECT 1 FROM product_variants v WHERE v.product_id=p.id AND v.active=true AND v.price_ars_cents >= "+next(f.MinPrice)+")")
	}
	if f.MaxPrice > 0 {
		where = append(where, "EXISTS (SELECT 1 FROM product_variants v WHERE v.product_id=p.id AND v.active=true AND v.price_ars_cents <= "+next(f.MaxPrice)+")")
	}
	countArgs := append([]any(nil), args...)
	var total int
	if err := r.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM products p WHERE `+strings.Join(where, " AND "), countArgs...).Scan(&total); err != nil {
		return ListResult{}, err
	}
	args = append(args, limit, f.Offset)
	query := `
SELECT p.id, p.slug, p.name, COALESCE(p.tagline,''), COALESCE(p.description,''), COALESCE(p.scent_family,''), COALESCE(p.gender_tag,''), COALESCE(p.concentration,''),
       COALESCE(p.top_notes, '{}'), COALESCE(p.heart_notes, '{}'), COALESCE(p.base_notes, '{}'), p.featured, p.active, p.display_order, p.created_at, p.updated_at,
       COALESCE((SELECT MIN(price_ars_cents) FROM product_variants v WHERE v.product_id=p.id AND v.active=true), 0),
       COALESCE((SELECT SUM(stock) FROM product_variants v WHERE v.product_id=p.id AND v.active=true), 0)
FROM products p
WHERE ` + strings.Join(where, " AND ") + `
ORDER BY p.display_order ASC, p.created_at DESC
LIMIT $` + fmt.Sprint(len(args)-1) + ` OFFSET $` + fmt.Sprint(len(args))
	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return ListResult{}, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		item, err := scanProduct(rows)
		if err != nil {
			return ListResult{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return ListResult{}, err
	}
	for i := range items {
		if err := r.hydrate(ctx, &items[i]); err != nil {
			return ListResult{}, err
		}
	}
	return ListResult{Items: items, Total: total, Limit: limit, Offset: f.Offset}, nil
}

func (r Repository) BySlug(ctx context.Context, slug string) (*Product, error) {
	if r.Pool == nil {
		return nil, pgx.ErrNoRows
	}
	row := r.Pool.QueryRow(ctx, `
SELECT p.id, p.slug, p.name, COALESCE(p.tagline,''), COALESCE(p.description,''), COALESCE(p.scent_family,''), COALESCE(p.gender_tag,''), COALESCE(p.concentration,''),
       COALESCE(p.top_notes, '{}'), COALESCE(p.heart_notes, '{}'), COALESCE(p.base_notes, '{}'), p.featured, p.active, p.display_order, p.created_at, p.updated_at,
       COALESCE((SELECT MIN(price_ars_cents) FROM product_variants v WHERE v.product_id=p.id AND v.active=true), 0),
       COALESCE((SELECT SUM(stock) FROM product_variants v WHERE v.product_id=p.id AND v.active=true), 0)
FROM products p WHERE p.slug=$1 AND p.active=true`, slug)
	p, err := scanProduct(row)
	if err != nil {
		return nil, err
	}
	return &p, r.hydrate(ctx, &p)
}

func (r Repository) Search(ctx context.Context, q string) ([]SearchResult, error) {
	if r.Pool == nil || strings.TrimSpace(q) == "" {
		return []SearchResult{}, nil
	}
	rows, err := r.Pool.Query(ctx, `
SELECT p.id, p.name, p.slug,
       COALESCE((SELECT url FROM product_images i WHERE i.product_id=p.id ORDER BY i.is_primary DESC, i.sort_order ASC LIMIT 1), ''),
       COALESCE((SELECT MIN(price_ars_cents) FROM product_variants v WHERE v.product_id=p.id AND v.active=true), 0)
FROM products p
WHERE p.active=true AND (p.name ILIKE $1 OR p.tagline ILIKE $1)
ORDER BY p.display_order ASC, p.name ASC
LIMIT 8`, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SearchResult{}
	for rows.Next() {
		var item SearchResult
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug, &item.PrimaryImage, &item.MinPrice); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r Repository) hydrate(ctx context.Context, p *Product) error {
	variants, err := r.variants(ctx, p.ID)
	if err != nil {
		return err
	}
	images, err := r.images(ctx, p.ID)
	if err != nil {
		return err
	}
	p.Variants = variants
	p.Images = images
	return nil
}

func (r Repository) variants(ctx context.Context, productID string) ([]Variant, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id, product_id, size_ml, price_ars_cents, stock, COALESCE(sku,''), active FROM product_variants WHERE product_id=$1 AND active=true ORDER BY size_ml`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Variant{}
	for rows.Next() {
		var v Variant
		if err := rows.Scan(&v.ID, &v.ProductID, &v.SizeML, &v.PriceARSCents, &v.Stock, &v.SKU, &v.Active); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r Repository) images(ctx context.Context, productID string) ([]ProductImage, error) {
	rows, err := r.Pool.Query(ctx, `SELECT id, product_id, url, COALESCE(alt_text,''), is_primary, sort_order FROM product_images WHERE product_id=$1 ORDER BY is_primary DESC, sort_order ASC`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ProductImage{}
	for rows.Next() {
		var img ProductImage
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.AltText, &img.IsPrimary, &img.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, img)
	}
	return out, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanProduct(row scanner) (Product, error) {
	var p Product
	err := row.Scan(&p.ID, &p.Slug, &p.Name, &p.Tagline, &p.Description, &p.ScentFamily, &p.GenderTag, &p.Concentration, &p.TopNotes, &p.HeartNotes, &p.BaseNotes, &p.Featured, &p.Active, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt, &p.MinPriceCents, &p.TotalStock)
	return p, err
}
