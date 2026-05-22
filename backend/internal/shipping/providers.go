package shipping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	correoArgentinoRatesURL = "https://api.correoargentino.com.ar/micorreo/v1/rates"
	andreaniLoginURL        = "https://apis.andreani.com/login"
	andreaniQuoteURL        = "https://apis.andreani.com/v1/cotizaciones"
)

type LocalPickupProvider struct{}

func (LocalPickupProvider) Quote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error) {
	return []QuoteOption{{
		ID:               "local-pickup",
		CarrierName:      "ELIXIR",
		ServiceName:      "Retiro en local",
		PriceCents:       0,
		EstimatedDaysMin: 0,
		EstimatedDaysMax: 0,
	}}, nil
}

func (LocalPickupProvider) CreateShipment(ctx context.Context, req ShipmentRequest) (ShipmentResponse, error) {
	return ShipmentResponse{CarrierName: "ELIXIR"}, nil
}

func (LocalPickupProvider) Track(ctx context.Context, trackingNumber string) (TrackingStatus, error) {
	return TrackingStatus{CarrierName: "ELIXIR", TrackingNumber: trackingNumber, Status: "pickup"}, nil
}

type CorreoArgentinoProvider struct {
	HTTP             *http.Client
	APIKey           string
	ClientID         string
	OriginPostalCode string
}

func (p CorreoArgentinoProvider) Quote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error) {
	if strings.TrimSpace(p.APIKey) == "" || strings.TrimSpace(p.ClientID) == "" {
		return nil, ErrProviderUnavailable
	}
	body, _ := json.Marshal(map[string]any{
		"origin_postal_code":      p.OriginPostalCode,
		"destination_postal_code": req.DestinationPostalCode,
		"weight_grams":            req.WeightGrams,
		"dimensions":              req.Dimensions,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, correoArgentinoRatesURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", p.APIKey)
	httpReq.Header.Set("X-Client-ID", p.ClientID)
	resp, err := httpClient(p.HTTP).Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("correo argentino quote failed: %s", resp.Status)
	}
	var raw map[string]any
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	return quoteOptionsFromResponse(raw, "Correo Argentino", "correo-argentino"), nil
}

func (p CorreoArgentinoProvider) CreateShipment(ctx context.Context, req ShipmentRequest) (ShipmentResponse, error) {
	return ShipmentResponse{}, ErrProviderUnavailable
}

func (p CorreoArgentinoProvider) Track(ctx context.Context, trackingNumber string) (TrackingStatus, error) {
	return TrackingStatus{}, ErrProviderUnavailable
}

type AndreaniProvider struct {
	HTTP             *http.Client
	User             string
	Password         string
	ClientID         string
	OriginPostalCode string
	mu               sync.Mutex
	token            string
	tokenExpiresAt   time.Time
}

func (p *AndreaniProvider) Quote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error) {
	if strings.TrimSpace(p.User) == "" || strings.TrimSpace(p.Password) == "" || strings.TrimSpace(p.ClientID) == "" {
		return nil, ErrProviderUnavailable
	}
	token, err := p.accessToken(ctx)
	if err != nil {
		return nil, err
	}
	body, _ := json.Marshal(map[string]any{
		"codigoPostalOrigen":  p.OriginPostalCode,
		"codigoPostalDestino": req.DestinationPostalCode,
		"peso":                req.WeightGrams,
		"alto":                req.Dimensions.HeightCM,
		"ancho":               req.Dimensions.WidthCM,
		"largo":               req.Dimensions.LengthCM,
		"cliente":             p.ClientID,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, andreaniQuoteURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("X-Cliente", p.ClientID)
	resp, err := httpClient(p.HTTP).Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("andreani quote failed: %s", resp.Status)
	}
	var raw map[string]any
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}
	return quoteOptionsFromResponse(raw, "Andreani", "andreani"), nil
}

func (p *AndreaniProvider) CreateShipment(ctx context.Context, req ShipmentRequest) (ShipmentResponse, error) {
	return ShipmentResponse{}, ErrProviderUnavailable
}

func (p *AndreaniProvider) Track(ctx context.Context, trackingNumber string) (TrackingStatus, error) {
	return TrackingStatus{}, ErrProviderUnavailable
}

func (p *AndreaniProvider) accessToken(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.token != "" && time.Now().Before(p.tokenExpiresAt) {
		return p.token, nil
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, andreaniLoginURL, nil)
	if err != nil {
		return "", err
	}
	httpReq.SetBasicAuth(p.User, p.Password)
	httpReq.Header.Set("X-Cliente", p.ClientID)
	resp, err := httpClient(p.HTTP).Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("andreani login failed: %s", resp.Status)
	}
	var data struct {
		AccessToken string `json:"access_token"`
		Token       string `json:"token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	if err := dec.Decode(&data); err != nil {
		return "", err
	}
	token := data.AccessToken
	if token == "" {
		token = data.Token
	}
	if token == "" {
		return "", ErrProviderUnavailable
	}
	ttl := time.Duration(data.ExpiresIn) * time.Second
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	p.token = token
	p.tokenExpiresAt = time.Now().Add(ttl - time.Minute)
	return p.token, nil
}

func httpClient(client *http.Client) *http.Client {
	if client != nil {
		return client
	}
	return &http.Client{Timeout: 8 * time.Second}
}

func quoteOptionsFromResponse(raw map[string]any, carrierName, idPrefix string) []QuoteOption {
	values := firstArray(raw, "rates", "options", "data", "cotizaciones")
	options := make([]QuoteOption, 0, len(values))
	for i, value := range values {
		item, ok := value.(map[string]any)
		if !ok {
			continue
		}
		service := firstString(item, "service_name", "serviceName", "servicio", "descripcion", "name")
		if service == "" {
			service = "Envío a domicilio"
		}
		price := firstCents(item, "price_cents", "priceCents", "importeConIva", "importe", "price", "total")
		if price < 0 {
			continue
		}
		id := strings.ToLower(strings.ReplaceAll(service, " ", "-"))
		options = append(options, QuoteOption{
			ID:               fmt.Sprintf("%s-%d-%s", idPrefix, i+1, id),
			CarrierName:      firstNonEmpty(firstString(item, "carrier_name", "carrierName", "carrier"), carrierName),
			ServiceName:      service,
			PriceCents:       price,
			EstimatedDaysMin: firstInt(item, "estimated_days_min", "estimatedDaysMin", "diasMin", "plazoMinimo"),
			EstimatedDaysMax: firstInt(item, "estimated_days_max", "estimatedDaysMax", "diasMax", "plazoMaximo"),
		})
	}
	return options
}

func firstArray(raw map[string]any, keys ...string) []any {
	for _, key := range keys {
		if values, ok := raw[key].([]any); ok {
			return values
		}
	}
	return nil
}

func firstString(raw map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := raw[key].(string); ok {
			return value
		}
	}
	return ""
}

func firstInt(raw map[string]any, keys ...string) int {
	for _, key := range keys {
		if value, ok := intFromValue(raw[key]); ok {
			return value
		}
	}
	return 0
}

func firstCents(raw map[string]any, keys ...string) int64 {
	for _, key := range keys {
		value := raw[key]
		if strings.Contains(strings.ToLower(key), "cent") {
			if cents, ok := int64FromValue(value); ok {
				return cents
			}
			continue
		}
		if cents, ok := centsFromARSValue(value); ok {
			return cents
		}
	}
	return -1
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func intFromValue(value any) (int, bool) {
	n, ok := int64FromValue(value)
	return int(n), ok
}

func int64FromValue(value any) (int64, bool) {
	switch v := value.(type) {
	case json.Number:
		n, err := strconv.ParseInt(v.String(), 10, 64)
		if err == nil {
			return n, true
		}
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err == nil {
			return n, true
		}
	case int64:
		return v, true
	case int:
		return int64(v), true
	}
	return 0, false
}

func centsFromARSValue(value any) (int64, bool) {
	switch v := value.(type) {
	case json.Number:
		return decimalToCents(v.String())
	case string:
		return decimalToCents(v)
	}
	if n, ok := int64FromValue(value); ok {
		return n * 100, true
	}
	return 0, false
}

func decimalToCents(value string) (int64, bool) {
	clean := strings.ReplaceAll(strings.TrimSpace(value), ",", ".")
	if clean == "" {
		return 0, false
	}
	whole, frac, hasFrac := strings.Cut(clean, ".")
	wholeCents, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, false
	}
	cents := wholeCents * 100
	if hasFrac {
		if len(frac) > 2 {
			frac = frac[:2]
		}
		for len(frac) < 2 {
			frac += "0"
		}
		fracCents, err := strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return 0, false
		}
		cents += fracCents
	}
	return cents, true
}
