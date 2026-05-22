package orders

import "time"

type CartItemRequest struct {
	VariantID            string `json:"variant_id"`
	Quantity             int    `json:"quantity"`
	ClientUnitPriceCents int64  `json:"unit_price_ars_cents,omitempty"`
}

type CartValidationRequest struct {
	Items []CartItemRequest `json:"items"`
}

type ValidatedItem struct {
	ProductID         string `json:"product_id"`
	VariantID         string `json:"variant_id"`
	ProductName       string `json:"product_name"`
	Slug              string `json:"slug"`
	PrimaryImage      string `json:"primary_image,omitempty"`
	SizeML            int    `json:"size_ml"`
	Quantity          int    `json:"quantity"`
	UnitPriceARSCents int64  `json:"unit_price_ars_cents"`
	SubtotalARSCents  int64  `json:"subtotal_ars_cents"`
	AvailableStock    int    `json:"available_stock"`
	CorrectedPrice    bool   `json:"corrected_price"`
	WeightGrams       int    `json:"weight_grams"`
}

type CartValidationResult struct {
	Valid            bool            `json:"valid"`
	Items            []ValidatedItem `json:"items"`
	SubtotalARSCents int64           `json:"subtotal_ars_cents"`
	Errors           []string        `json:"errors"`
}

type CreateOrderRequest struct {
	Items                []CartItemRequest `json:"items"`
	CustomerName         string            `json:"customer_name"`
	CustomerEmail        string            `json:"customer_email"`
	CustomerPhone        string            `json:"customer_phone"`
	ShippingAddress      map[string]any    `json:"shipping_address"`
	ShippingCostARSCents int64             `json:"shipping_cost_ars_cents"`
	DiscountCode         string            `json:"discount_code"`
	Notes                string            `json:"notes"`
}

type Order struct {
	ID                   string          `json:"id"`
	ExternalReference    string          `json:"external_reference"`
	Status               string          `json:"status"`
	CustomerName         string          `json:"customer_name"`
	CustomerEmail        string          `json:"customer_email"`
	CustomerPhone        string          `json:"customer_phone,omitempty"`
	ShippingAddress      map[string]any  `json:"shipping_address,omitempty"`
	ShippingCostARSCents int64           `json:"shipping_cost_ars_cents"`
	SubtotalARSCents     int64           `json:"subtotal_ars_cents"`
	TotalARSCents        int64           `json:"total_ars_cents"`
	DiscountCode         string          `json:"discount_code,omitempty"`
	DiscountARSCents     int64           `json:"discount_ars_cents"`
	Currency             string          `json:"currency"`
	TrackingNumber       string          `json:"tracking_number,omitempty"`
	Notes                string          `json:"notes,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	Items                []ValidatedItem `json:"items"`
}
