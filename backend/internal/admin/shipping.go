package admin

import (
	"errors"
	"net/http"
	"strings"

	"elixir/backend/internal/audit"
	"elixir/backend/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type shippingZonePayload struct {
	ID               string   `json:"id,omitempty"`
	ZoneName         string   `json:"zone_name"`
	ProvinceCodes    []string `json:"province_codes"`
	BaseCostCents    int64    `json:"base_cost_cents"`
	PerKGCents       int64    `json:"per_kg_cents"`
	EstimatedDaysMin int      `json:"estimated_days_min"`
	EstimatedDaysMax int      `json:"estimated_days_max"`
	Active           bool     `json:"active"`
}

func (h Handler) AdminShippingZones(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		if r.Method == http.MethodGet {
			httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
			return
		}
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if r.Method == http.MethodGet {
		rows, err := h.Pool.Query(r.Context(), `SELECT id, zone_name, COALESCE(province_codes,'{}'), base_cost_cents, per_kg_cents, estimated_days_min, estimated_days_max, active FROM shipping_zones ORDER BY active DESC, base_cost_cents ASC, zone_name ASC`)
		if err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron listar zonas")
			return
		}
		defer rows.Close()
		items := []shippingZonePayload{}
		for rows.Next() {
			var item shippingZonePayload
			if err := rows.Scan(&item.ID, &item.ZoneName, &item.ProvinceCodes, &item.BaseCostCents, &item.PerKGCents, &item.EstimatedDaysMin, &item.EstimatedDaysMax, &item.Active); err != nil {
				httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer zonas")
				return
			}
			items = append(items, item)
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}
	var req shippingZonePayload
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if err := normalizeShippingZonePayload(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	var id string
	err := h.Pool.QueryRow(r.Context(), `
		INSERT INTO shipping_zones (zone_name, province_codes, base_cost_cents, per_kg_cents, estimated_days_min, estimated_days_max, active)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id`,
		req.ZoneName, req.ProvinceCodes, req.BaseCostCents, req.PerKGCents, req.EstimatedDaysMin, req.EstimatedDaysMax, req.Active,
	).Scan(&id)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear la zona de envío")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "shipping_zone.create", ResourceID: id, Metadata: map[string]any{"name": req.ZoneName}})
	httpx.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (h Handler) AdminShippingZoneByID(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	id := chi.URLParam(r, "id")
	if r.Method == http.MethodDelete {
		tag, err := h.Pool.Exec(r.Context(), `UPDATE shipping_zones SET active=false WHERE id=$1`, id)
		if err != nil || tag.RowsAffected() == 0 {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo desactivar la zona")
			return
		}
		audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "shipping_zone.deactivate", ResourceID: id})
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	var req shippingZonePayload
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if err := normalizeShippingZonePayload(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	tag, err := h.Pool.Exec(r.Context(), `
		UPDATE shipping_zones
		SET zone_name=$1, province_codes=$2, base_cost_cents=$3, per_kg_cents=$4, estimated_days_min=$5, estimated_days_max=$6, active=$7
		WHERE id=$8`,
		req.ZoneName, req.ProvinceCodes, req.BaseCostCents, req.PerKGCents, req.EstimatedDaysMin, req.EstimatedDaysMax, req.Active, id,
	)
	if err != nil || tag.RowsAffected() == 0 {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar la zona")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "shipping_zone.update", ResourceID: id, Metadata: map[string]any{"name": req.ZoneName, "active": req.Active}})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func normalizeShippingZonePayload(req *shippingZonePayload) error {
	req.ZoneName = strings.TrimSpace(req.ZoneName)
	if req.ZoneName == "" || len(req.ZoneName) > 80 {
		return errors.New("la zona necesita un nombre claro")
	}
	if req.BaseCostCents < 0 || req.PerKGCents < 0 {
		return errors.New("los costos de envío no pueden ser negativos")
	}
	if req.EstimatedDaysMin < 0 || req.EstimatedDaysMax < 0 || req.EstimatedDaysMin > req.EstimatedDaysMax || req.EstimatedDaysMax > 60 {
		return errors.New("revisá los días estimados de entrega")
	}
	req.ProvinceCodes = cleanProvinceCodes(req.ProvinceCodes)
	if len(req.ProvinceCodes) == 0 {
		req.ProvinceCodes = []string{"AR"}
	}
	return nil
}

func cleanProvinceCodes(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			code := strings.ToUpper(strings.TrimSpace(part))
			code = strings.ReplaceAll(code, " ", "_")
			if code == "" || seen[code] {
				continue
			}
			seen[code] = true
			out = append(out, code)
			if len(out) == 24 {
				return out
			}
		}
	}
	return out
}
