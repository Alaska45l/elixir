package payments

import (
	"context"
	"testing"
)

type fakePaymentRepo struct {
	exists   bool
	recorded bool
}

func (f *fakePaymentRepo) OrderForPayment(context.Context, string) (OrderForPayment, error) {
	return OrderForPayment{}, nil
}
func (f *fakePaymentRepo) HasPaymentEvent(context.Context, string) (bool, error) {
	return f.exists, nil
}
func (f *fakePaymentRepo) RecordPaymentAndUpdateOrder(context.Context, PaymentDetails) error {
	f.recorded = true
	return nil
}

func TestWebhookIdempotencySkipsDuplicate(t *testing.T) {
	repo := &fakePaymentRepo{exists: true}
	svc := Service{Repo: repo, MP: MercadoPagoClient{}}
	if err := svc.ProcessWebhook(context.Background(), "123"); err != nil {
		t.Fatal(err)
	}
	if repo.recorded {
		t.Fatal("duplicate payment event should not be recorded")
	}
}
