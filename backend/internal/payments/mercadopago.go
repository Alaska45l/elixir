package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type MercadoPagoClient struct {
	AccessToken string
	BackendURL  string
	FrontendURL string
	HTTPClient  *http.Client
}

func (c MercadoPagoClient) CreatePreference(ctx context.Context, order OrderForPayment) (PreferenceResponse, error) {
	if c.AccessToken == "" {
		return PreferenceResponse{
			PreferenceID: "local-" + order.ExternalReference,
			InitPoint:    c.FrontendURL + "/orden/" + order.ExternalReference + "?status=pending",
		}, nil
	}
	items := make([]map[string]any, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, map[string]any{
			"title":       item.Title,
			"quantity":    item.Quantity,
			"currency_id": "ARS",
			"unit_price":  float64(item.UnitPriceARSCents) / 100,
		})
	}
	payload := map[string]any{
		"items":              items,
		"external_reference": order.ExternalReference,
		"payer":              map[string]any{"name": order.CustomerName, "email": order.CustomerEmail},
		"back_urls": map[string]any{
			"success": c.FrontendURL + "/orden/" + order.ExternalReference + "?status=success",
			"failure": c.FrontendURL + "/carrito?status=failed",
			"pending": c.FrontendURL + "/orden/" + order.ExternalReference + "?status=pending",
		},
		"notification_url": c.BackendURL + "/api/payments/mercadopago/webhook",
		"binary_mode":      false,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return PreferenceResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.mercadopago.com/checkout/preferences", bytes.NewReader(body))
	if err != nil {
		return PreferenceResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	res, err := client.Do(req)
	if err != nil {
		return PreferenceResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return PreferenceResponse{}, fmt.Errorf("mercadopago preference status %d", res.StatusCode)
	}
	var out struct {
		ID        string `json:"id"`
		InitPoint string `json:"init_point"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return PreferenceResponse{}, err
	}
	return PreferenceResponse{PreferenceID: out.ID, InitPoint: out.InitPoint}, nil
}

func (c MercadoPagoClient) FetchPayment(ctx context.Context, id string) (PaymentDetails, error) {
	if id == "" {
		return PaymentDetails{}, errors.New("missing payment id")
	}
	if c.AccessToken == "" {
		return PaymentDetails{ID: id, Status: "in_process", Raw: map[string]any{"id": id, "local": true}}, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.mercadopago.com/v1/payments/"+id, nil)
	if err != nil {
		return PaymentDetails{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	res, err := client.Do(req)
	if err != nil {
		return PaymentDetails{}, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return PaymentDetails{}, fmt.Errorf("mercadopago payment status %d", res.StatusCode)
	}
	var raw map[string]any
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return PaymentDetails{}, err
	}
	return PaymentDetails{
		ID:                stringValue(raw["id"]),
		Status:            stringValue(raw["status"]),
		StatusDetail:      stringValue(raw["status_detail"]),
		ExternalReference: stringValue(raw["external_reference"]),
		PreferenceID:      stringValue(raw["preference_id"]),
		Raw:               raw,
	}, nil
}

func stringValue(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%.0f", t)
	default:
		return ""
	}
}
