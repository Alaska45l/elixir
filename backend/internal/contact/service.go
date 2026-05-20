package contact

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Pool *pgxpool.Pool
}

type MessageRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AbandonedCartRequest struct {
	Email    string         `json:"email"`
	CartData map[string]any `json:"cart_data"`
}

func (s Service) SaveMessage(ctx context.Context, req MessageRequest) error {
	if s.Pool == nil {
		return nil
	}
	_, err := s.Pool.Exec(ctx, `INSERT INTO contact_messages (name,email,subject,message) VALUES ($1,$2,$3,$4)`, strings.TrimSpace(req.Name), strings.TrimSpace(req.Email), req.Subject, strings.TrimSpace(req.Message))
	return err
}

func (s Service) SaveAbandoned(ctx context.Context, req AbandonedCartRequest) error {
	if s.Pool == nil || strings.TrimSpace(req.Email) == "" {
		return nil
	}
	body, _ := json.Marshal(req.CartData)
	_, err := s.Pool.Exec(ctx, `INSERT INTO abandoned_carts (email, cart_data) VALUES ($1,$2)`, strings.TrimSpace(req.Email), body)
	return err
}
