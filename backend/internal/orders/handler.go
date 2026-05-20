package orders

import (
	"errors"
	"net/http"

	"elixir/backend/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	Service Service
}

func (h Handler) ValidateCart(w http.ResponseWriter, r *http.Request) {
	var req CartValidationRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	res, err := h.Service.ValidateCart(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo validar el carrito")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	order, err := h.Service.CreateOrder(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, order)
}

func (h Handler) Status(w http.ResponseWriter, r *http.Request) {
	order, err := h.Service.Repo.ByExternalReference(r.Context(), chi.URLParam(r, "external_reference"))
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, r, http.StatusNotFound, "orden no encontrada")
		return
	}
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo obtener la orden")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, order)
}
