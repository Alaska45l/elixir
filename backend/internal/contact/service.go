package contact

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Pool *pgxpool.Pool
}

var errInvalidContactInput = errors.New("contact input inválido")

type MessageRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Website string `json:"website,omitempty"`
}

type AbandonedCartRequest struct {
	Email    string         `json:"email"`
	CartData map[string]any `json:"cart_data"`
}

func (s Service) SaveMessage(ctx context.Context, req MessageRequest) error {
	if strings.TrimSpace(req.Website) != "" {
		return nil
	}
	if err := normalizeMessage(&req); err != nil {
		return err
	}
	if s.Pool == nil {
		return nil
	}
	_, err := s.Pool.Exec(ctx, `INSERT INTO contact_messages (name,email,subject,message) VALUES ($1,$2,$3,$4)`, req.Name, req.Email, req.Subject, req.Message)
	return err
}

func (s Service) SaveAbandoned(ctx context.Context, req AbandonedCartRequest) error {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" {
		return nil
	}
	if len(req.Email) > 254 {
		return fmt.Errorf("%w: email inválido", errInvalidContactInput)
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("%w: email inválido", errInvalidContactInput)
	}
	body, err := json.Marshal(req.CartData)
	if err != nil {
		return err
	}
	if len(body) > 64*1024 {
		return fmt.Errorf("%w: carrito demasiado grande", errInvalidContactInput)
	}
	if s.Pool == nil {
		return nil
	}
	_, err = s.Pool.Exec(ctx, `INSERT INTO abandoned_carts (email, cart_data) VALUES ($1,$2)`, req.Email, body)
	return err
}

func normalizeMessage(req *MessageRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Subject = strings.TrimSpace(req.Subject)
	req.Message = strings.TrimSpace(req.Message)
	if strings.TrimSpace(req.Website) != "" {
		return nil
	}
	if len(req.Name) < 2 || len(req.Name) > 120 {
		return fmt.Errorf("%w: nombre inválido", errInvalidContactInput)
	}
	if len(req.Email) > 254 {
		return fmt.Errorf("%w: email inválido", errInvalidContactInput)
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("%w: email inválido", errInvalidContactInput)
	}
	if len(req.Subject) > 160 {
		return fmt.Errorf("%w: asunto demasiado largo", errInvalidContactInput)
	}
	if len(req.Message) < 10 || len(req.Message) > 4000 {
		return fmt.Errorf("%w: mensaje inválido", errInvalidContactInput)
	}
	return nil
}
