package payments

import (
	"context"
)

type Repository interface {
	OrderForPayment(ctx context.Context, ref string) (OrderForPayment, error)
	HasPaymentEvent(ctx context.Context, paymentID string) (bool, error)
	RecordPaymentAndUpdateOrder(ctx context.Context, payment PaymentDetails) error
}

type Service struct {
	Repo Repository
	MP   MercadoPagoClient
}

func (s Service) CreatePreference(ctx context.Context, ref string) (PreferenceResponse, error) {
	order, err := s.Repo.OrderForPayment(ctx, ref)
	if err != nil {
		return PreferenceResponse{}, err
	}
	return s.MP.CreatePreference(ctx, order)
}

func (s Service) ProcessWebhook(ctx context.Context, paymentID string) error {
	exists, err := s.Repo.HasPaymentEvent(ctx, paymentID)
	if err != nil || exists {
		return err
	}
	payment, err := s.MP.FetchPayment(ctx, paymentID)
	if err != nil {
		return err
	}
	if payment.ID == "" {
		payment.ID = paymentID
	}
	return s.Repo.RecordPaymentAndUpdateOrder(ctx, payment)
}
