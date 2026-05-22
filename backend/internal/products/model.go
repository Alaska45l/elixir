package products

import "time"

type Product struct {
	ID            string         `json:"id"`
	Slug          string         `json:"slug"`
	Name          string         `json:"name"`
	Tagline       string         `json:"tagline,omitempty"`
	Description   string         `json:"description,omitempty"`
	ScentFamily   string         `json:"scent_family,omitempty"`
	GenderTag     string         `json:"gender_tag,omitempty"`
	Concentration string         `json:"concentration,omitempty"`
	TopNotes      []string       `json:"top_notes"`
	HeartNotes    []string       `json:"heart_notes"`
	BaseNotes     []string       `json:"base_notes"`
	Featured      bool           `json:"featured"`
	Active        bool           `json:"active"`
	DisplayOrder  int            `json:"display_order"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Variants      []Variant      `json:"variants"`
	Images        []ProductImage `json:"images"`
	MinPriceCents int64          `json:"min_price_ars_cents"`
	TotalStock    int            `json:"total_stock"`
}

type Variant struct {
	ID            string `json:"id"`
	ProductID     string `json:"product_id"`
	SizeML        int    `json:"size_ml"`
	PriceARSCents int64  `json:"price_ars_cents"`
	Stock         int    `json:"stock"`
	SKU           string `json:"sku,omitempty"`
	Active        bool   `json:"active"`
	WeightGrams   int    `json:"weight_grams"`
}

type ProductImage struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	URL       string `json:"url"`
	AltText   string `json:"alt_text,omitempty"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type ListFilters struct {
	Featured       *bool
	Families       []string
	Genders        []string
	Concentrations []string
	MinPrice       int64
	MaxPrice       int64
	InStock        bool
	Search         string
	Limit          int
	Offset         int
}

type ListResult struct {
	Items  []Product `json:"items"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

type SearchResult struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	PrimaryImage string `json:"primary_image,omitempty"`
	MinPrice     int64  `json:"min_price_ars_cents"`
}
