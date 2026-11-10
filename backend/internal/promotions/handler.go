package promotions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ProductGetter interface {
	GetPriceAndName(id int64) (string, float64, error)
}

type Handler struct {
	repo          *Repository
	productGetter ProductGetter
}

func NewHandler(repo *Repository, productGetter ProductGetter) *Handler {
	return &Handler{repo: repo, productGetter: productGetter}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/promociones", h.list)
	mux.HandleFunc("GET /api/promociones/{id}", h.get)
	mux.HandleFunc("POST /api/promociones", h.create)
	mux.HandleFunc("PUT /api/promociones/{id}", h.update)
	mux.HandleFunc("DELETE /api/promociones/{id}", h.delete)
	mux.HandleFunc("POST /api/promociones/aplicar", h.apply)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudieron listar las promociones: "+err.Error())
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
		respondError(w, http.StatusNotFound, "Promoción no encontrada")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo consultar la promoción: "+err.Error())
		return
	}
	respond(w, http.StatusOK, p)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in PromotionInput
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
		respondError(w, http.StatusInternalServerError, "No se pudo registrar la promoción: "+err.Error())
		return
	}
	respond(w, http.StatusCreated, p)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var in PromotionInput
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
		respondError(w, http.StatusNotFound, "Promoción no encontrada")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo actualizar la promoción: "+err.Error())
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
		respondError(w, http.StatusNotFound, "Promoción no encontrada")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "No se pudo eliminar la promoción: "+err.Error())
		return
	}
	respond(w, http.StatusOK, map[string]string{"mensaje": "Promoción eliminada exitosamente"})
}

func (h *Handler) apply(w http.ResponseWriter, r *http.Request) {
	var req ApplyPromotionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	prodName := "Carrito de Compras"
	basePrice := req.CustomAmount

	if req.ProductID > 0 && h.productGetter != nil {
		name, price, err := h.productGetter.GetPriceAndName(req.ProductID)
		if err == nil {
			prodName = name
			basePrice = price
		}
	}

	if basePrice <= 0 {
		basePrice = 100.0 // Monto base referencial por defecto
	}

	// 1. Instanciar el Componente Concreto Base
	var currentComponent PriceComponent = NewBasePrice(prodName, basePrice)
	decoratorCount := 0

	// 2. Aplicar decoradores para cada promoción seleccionada por ID
	for _, promoID := range req.PromotionIDs {
		promo, err := h.repo.Get(r.Context(), promoID)
		if err != nil || !promo.Active {
			continue
		}
		decoratorCount++
		switch promo.Type {
		case "porcentaje":
			currentComponent = NewPercentageDecorator(currentComponent, promo.Name, promo.Value)
		case "monto_fijo":
			currentComponent = NewFixedAmountDecorator(currentComponent, promo.Name, promo.Value)
		case "cupon":
			currentComponent = NewCouponDecorator(currentComponent, promo.Code, promo.Name, promo.Type, promo.Value, promo.MinAmount)
		}
	}

	// 3. Aplicar cupón por código si fue proporcionado
	if strings.TrimSpace(req.CouponCode) != "" {
		coupon, err := h.repo.GetByCode(r.Context(), req.CouponCode)
		if err == nil && coupon.Active {
			decoratorCount++
			currentComponent = NewCouponDecorator(currentComponent, coupon.Code, coupon.Name, coupon.Type, coupon.Value, coupon.MinAmount)
		}
	}

	finalPrice := currentComponent.GetFinalPrice()
	totalDiscount := basePrice - finalPrice
	if totalDiscount < 0 {
		totalDiscount = 0
	}

	res := ApplyPromotionResponse{
		ProductID:       req.ProductID,
		ProductName:     prodName,
		BasePrice:       basePrice,
		TotalDiscount:   totalDiscount,
		FinalPrice:      finalPrice,
		DecoratorsCount: decoratorCount,
		Steps:           currentComponent.GetBreakdown(),
	}

	respond(w, http.StatusOK, res)
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
