package admin

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"elixir/backend/internal/audit"
	"elixir/backend/internal/httpx"
	"elixir/backend/internal/media"
)

const maxUploadBytes = 10 << 20

func (h Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if !h.Media.IsConfigured() {
		httpx.Error(w, r, http.StatusServiceUnavailable, "almacenamiento de imágenes no configurado")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		status := http.StatusBadRequest
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			status = http.StatusRequestEntityTooLarge
		}
		httpx.Error(w, r, status, "archivo inválido o demasiado grande (máximo 10 MB)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "campo 'file' requerido")
		return
	}
	defer file.Close()
	if err := validateUploadHeader(header.Filename, header.Header.Get("Content-Type")); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	folder := r.FormValue("folder")
	if folder == "" {
		folder = "products"
	}

	url, err := h.Media.ProcessAndUpload(r.Context(), file, folder)
	if err != nil {
		if errors.Is(err, media.ErrInvalidImage) {
			httpx.Error(w, r, http.StatusBadRequest, "archivo de imagen inválido")
			return
		}
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo procesar la imagen")
		return
	}

	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "media.upload", Metadata: map[string]any{"folder": folder}})
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"url": url})
}

func validateUploadHeader(filename, contentType string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return errors.New("la extensión debe ser jpg, jpeg, png o webp")
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return errors.New("content-type de imagen inválido")
	}
	switch strings.ToLower(mediaType) {
	case "image/jpeg", "image/png", "image/webp":
		return nil
	default:
		return errors.New("content-type de imagen no permitido")
	}
}
