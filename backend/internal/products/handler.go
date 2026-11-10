package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Handler struct{ repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{repo: repo} }

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/productos", h.list)
	mux.HandleFunc("GET /api/productos/{id}", h.get)
	mux.HandleFunc("GET /api/productos/{id}/composite", h.getComposite)
	mux.HandleFunc("POST /api/productos", h.create)
	mux.HandleFunc("PUT /api/productos/{id}", h.update)
	mux.HandleFunc("DELETE /api/productos/{id}", h.delete)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudieron consultar los productos: "+err.Error())
		return
	}
	respond(w, http.StatusOK, items)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	p, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo consultar el producto: "+err.Error())
		return
	}
	respond(w, http.StatusOK, p)
}

func (h *Handler) getComposite(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	p, err := h.repo.Get(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error al consultar: "+err.Error())
		return
	}

	var rawPrice float64
	for _, it := range p.Items {
		rawPrice += it.Subtotal
	}
	savings := rawPrice - p.Price
	if savings < 0 {
		savings = 0
	}

	res := CompositeDetailResponse{
		ID:              p.ID,
		Nombre:          p.Name,
		Tipo:            p.Type,
		PrecioFinal:     p.Price,
		PrecioBruto:     rawPrice,
		Ahorro:          savings,
		DescuentoPct:    p.DiscountPercent,
		StockDisponible: p.Stock,
		Hijos:           p.Items,
	}
	respond(w, http.StatusOK, res)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in ProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := ValidateInput(in); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err := h.repo.Create(r.Context(), in)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo crear el producto: "+err.Error())
		return
	}
	respond(w, http.StatusCreated, p)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in ProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := ValidateInput(in); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	p, err := h.repo.Update(r.Context(), id, in)
	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo actualizar el producto: "+err.Error())
		return
	}
	respond(w, http.StatusOK, p)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	err := h.repo.Delete(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "Producto no encontrado")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo eliminar el producto: "+err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"mensaje": "Producto eliminado exitosamente"})
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		respondError(w, http.StatusBadRequest, "ID requerido en la ruta")
		return 0, false
	}
	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "ID inválido")
		return 0, false
	}
	return id, true
}

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respond(w, status, map[string]string{"error": message})
}
