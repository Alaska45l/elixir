package payments

type PreferenceRequest struct {
	ExternalReference string `json:"external_reference"`
}

type PreferenceResponse struct {
	PreferenceID string `json:"preference_id"`
	InitPoint    string `json:"init_point"`
}

type PaymentDetails struct {
	ID                string         `json:"id"`
	Status            string         `json:"status"`
	StatusDetail      string         `json:"status_detail"`
	ExternalReference string         `json:"external_reference"`
	PreferenceID      string         `json:"preference_id"`
	Raw               map[string]any `json:"raw"`
}

type OrderForPayment struct {
	ID                string
	ExternalReference string
	CustomerEmail     string
	CustomerName      string
	TotalARSCents     int64
	Items             []PaymentItem
}

type PaymentItem struct {
	Title             string
	Quantity          int
	UnitPriceARSCents int64
	VariantID         string
}
