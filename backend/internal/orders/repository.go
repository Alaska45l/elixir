package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepository struct {
	Pool *pgxpool.Pool
}

func (r DBRepository) FetchVariant(ctx context.Context, variantID string) (ValidatedItem, error) {
	if r.Pool == nil {
		return ValidatedItem{}, pgx.ErrNoRows
	}
	row := r.Pool.QueryRow(ctx, `
SELECT p.id, v.id, p.name, p.slug,
       COALESCE((SELECT url FROM product_images i WHERE i.product_id=p.id ORDER BY i.is_primary DESC, i.sort_order ASC LIMIT 1), ''),
	       v.size_ml, v.price_ars_cents, v.stock, v.weight_grams
FROM product_variants v
JOIN products p ON p.id=v.product_id
WHERE v.id=$1 AND v.active=true AND p.active=true`, variantID)
	var item ValidatedItem
	err := row.Scan(&item.ProductID, &item.VariantID, &item.ProductName, &item.Slug, &item.PrimaryImage, &item.SizeML, &item.UnitPriceARSCents, &item.AvailableStock, &item.WeightGrams)
	return item, err
}

func (r DBRepository) Create(ctx context.Context, req CreateOrderRequest, validation CartValidationResult, discountCents int64) (*Order, error) {
	if r.Pool == nil {
		return nil, errors.New("database not configured")
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	ref := "ELX-" + uuid.NewString()
	total := CalculateTotal(validation.SubtotalARSCents, req.ShippingCostARSCents, discountCents)
	address, _ := json.Marshal(req.ShippingAddress)
	var order Order
	err = tx.QueryRow(ctx, `
INSERT INTO orders (external_reference, customer_name, customer_email, customer_phone, shipping_address, shipping_cost_ars_cents, subtotal_ars_cents, total_ars_cents, discount_code, discount_ars_cents, notes)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
RETURNING id, external_reference, status, customer_name, customer_email, COALESCE(customer_phone,''), shipping_cost_ars_cents, subtotal_ars_cents, total_ars_cents, COALESCE(discount_code,''), discount_ars_cents, currency, COALESCE(tracking_number,''), COALESCE(notes,''), created_at, updated_at`,
		ref, req.CustomerName, req.CustomerEmail, req.CustomerPhone, address, req.ShippingCostARSCents, validation.SubtotalARSCents, total, req.DiscountCode, discountCents, req.Notes,
	).Scan(&order.ID, &order.ExternalReference, &order.Status, &order.CustomerName, &order.CustomerEmail, &order.CustomerPhone, &order.ShippingCostARSCents, &order.SubtotalARSCents, &order.TotalARSCents, &order.DiscountCode, &order.DiscountARSCents, &order.Currency, &order.TrackingNumber, &order.Notes, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	for _, item := range validation.Items {
		tag, err := tx.Exec(ctx, `UPDATE product_variants SET stock = stock - $1 WHERE id=$2 AND active=true AND stock >= $1`, item.Quantity, item.VariantID)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() != 1 {
			return nil, fmt.Errorf("%s no tiene stock suficiente", item.ProductName)
		}
		_, err = tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, variant_id, product_name, size_ml, quantity, unit_price_ars_cents, subtotal_ars_cents) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			order.ID, item.ProductID, item.VariantID, item.ProductName, item.SizeML, item.Quantity, item.UnitPriceARSCents, item.SubtotalARSCents)
		if err != nil {
			return nil, err
		}
	}
	order.Items = validation.Items
	return &order, tx.Commit(ctx)
}

func (r DBRepository) ByExternalReference(ctx context.Context, ref string) (*Order, error) {
	if r.Pool == nil {
		return nil, pgx.ErrNoRows
	}
	var raw []byte
	row := r.Pool.QueryRow(ctx, `
SELECT id, external_reference, status, customer_name, customer_email, COALESCE(customer_phone,''), COALESCE(shipping_address,'{}'::jsonb), shipping_cost_ars_cents, subtotal_ars_cents, total_ars_cents, COALESCE(discount_code,''), discount_ars_cents, currency, COALESCE(tracking_number,''), COALESCE(notes,''), created_at, updated_at
FROM orders WHERE external_reference=$1`, ref)
	var o Order
	if err := row.Scan(&o.ID, &o.ExternalReference, &o.Status, &o.CustomerName, &o.CustomerEmail, &o.CustomerPhone, &raw, &o.ShippingCostARSCents, &o.SubtotalARSCents, &o.TotalARSCents, &o.DiscountCode, &o.DiscountARSCents, &o.Currency, &o.TrackingNumber, &o.Notes, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(raw, &o.ShippingAddress)
	return &o, nil
}
