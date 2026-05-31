package admin

import (
	"errors"
	"net/http"

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

	file, _, err := r.FormFile("file")
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "campo 'file' requerido")
		return
	}
	defer file.Close()

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

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"url": url})
}
