package shipping

import (
	"context"
	"errors"
)

var ErrProviderUnavailable = errors.New("shipping provider unavailable")

type Dimensions struct {
	LengthCM int `json:"length_cm"`
	WidthCM  int `json:"width_cm"`
	HeightCM int `json:"height_cm"`
}

type QuoteRequest struct {
	DestinationPostalCode string     `json:"destination_postal_code"`
	ProvinceCode          string     `json:"province_code"`
	WeightGrams           int        `json:"weight_grams"`
	Dimensions            Dimensions `json:"dimensions"`
}

type QuoteOption struct {
	ID               string `json:"id"`
	CarrierName      string `json:"carrier_name"`
	ServiceName      string `json:"service_name"`
	PriceCents       int64  `json:"price_cents"`
	EstimatedDaysMin int    `json:"estimated_days_min"`
	EstimatedDaysMax int    `json:"estimated_days_max"`
}

type ShipmentRequest struct {
	OrderID               string     `json:"order_id"`
	DestinationPostalCode string     `json:"destination_postal_code"`
	ProvinceCode          string     `json:"province_code"`
	WeightGrams           int        `json:"weight_grams"`
	Dimensions            Dimensions `json:"dimensions"`
}

type ShipmentResponse struct {
	CarrierName    string `json:"carrier_name"`
	TrackingNumber string `json:"tracking_number"`
	LabelURL       string `json:"label_url"`
}

type TrackingStatus struct {
	CarrierName    string `json:"carrier_name"`
	TrackingNumber string `json:"tracking_number"`
	Status         string `json:"status"`
	Description    string `json:"description"`
}

type ShippingProvider interface {
	Quote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error)
	CreateShipment(ctx context.Context, req ShipmentRequest) (ShipmentResponse, error)
	Track(ctx context.Context, trackingNumber string) (TrackingStatus, error)
}
