package payments

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBRepository struct {
	Pool *pgxpool.Pool
}

func (r DBRepository) OrderForPayment(ctx context.Context, ref string) (OrderForPayment, error) {
	var order OrderForPayment
	if r.Pool == nil {
		return order, pgx.ErrNoRows
	}
	err := r.Pool.QueryRow(ctx, `SELECT id, external_reference, customer_email, customer_name, total_ars_cents FROM orders WHERE external_reference=$1`, ref).
		Scan(&order.ID, &order.ExternalReference, &order.CustomerEmail, &order.CustomerName, &order.TotalARSCents)
	if err != nil {
		return order, err
	}
	rows, err := r.Pool.Query(ctx, `SELECT product_name, quantity, unit_price_ars_cents, variant_id FROM order_items WHERE order_id=$1`, order.ID)
	if err != nil {
		return order, err
	}
	defer rows.Close()
	for rows.Next() {
		var item PaymentItem
		if err := rows.Scan(&item.Title, &item.Quantity, &item.UnitPriceARSCents, &item.VariantID); err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}
	return order, rows.Err()
}

func (r DBRepository) HasPaymentEvent(ctx context.Context, paymentID string) (bool, error) {
	if r.Pool == nil {
		return false, nil
	}
	var exists bool
	err := r.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM payment_events WHERE mp_payment_id=$1)`, paymentID).Scan(&exists)
	return exists, err
}

func (r DBRepository) RecordPaymentAndUpdateOrder(ctx context.Context, payment PaymentDetails) error {
	if r.Pool == nil {
		return nil
	}
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var orderID *string
	if payment.ExternalReference != "" {
		var id string
		if err := tx.QueryRow(ctx, `SELECT id FROM orders WHERE external_reference=$1`, payment.ExternalReference).Scan(&id); err == nil {
			orderID = &id
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}
	raw, err := json.Marshal(payment.Raw)
	if err != nil {
		return err
	}
	tag, err := tx.Exec(ctx, `INSERT INTO payment_events (order_id, mp_payment_id, mp_preference_id, mp_status, mp_status_detail, raw_payload) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (mp_payment_id) DO NOTHING`,
		orderID, payment.ID, payment.PreferenceID, payment.Status, payment.StatusDetail, raw)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return tx.Commit(ctx)
	}
	if orderID != nil {
		status := mapStatus(payment.Status)
		_, err = tx.Exec(ctx, `
			UPDATE orders
			SET status=$1, updated_at=now()
			WHERE id=$2 AND (
				(status='pending' AND $1 IN ('paid','failed','cancelled')) OR
				(status='paid' AND $1='paid') OR
				(status=$1)
			)`, status, *orderID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func mapStatus(mp string) string {
	switch mp {
	case "approved":
		return "paid"
	case "rejected":
		return "failed"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}
