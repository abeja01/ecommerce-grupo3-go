package promotions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockProductGetter struct{}

func (m mockProductGetter) GetPriceAndName(id int64) (string, float64, error) {
	return "Laptop Pro 14", 3000.0, nil
}

func setupTestPromotionsRouter() *http.ServeMux {
	repo := NewRepository(nil)
	h := NewHandler(repo, mockProductGetter{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestPromotionsCRUD(t *testing.T) {
	mux := setupTestPromotionsRouter()

	// 1. Listar promociones (GET)
	req := httptest.NewRequest("GET", "/api/promociones", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en GET /api/promociones, obtenido %d", w.Code)
	}

	// 2. Crear promoción (POST)
	newPromo := PromotionInput{
		Code:        "DESC50",
		Name:        "Descuento 50 Soles",
		Description: "Bono especial",
		Type:        "monto_fijo",
		Value:       50.0,
		MinAmount:   200.0,
		Active:      true,
	}
	body, _ := json.Marshal(newPromo)
	req = httptest.NewRequest("POST", "/api/promociones", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Esperado 201 en POST /api/promociones, obtenido %d: %s", w.Code, w.Body.String())
	}

	// 3. Consultar promoción por ID (GET)
	req = httptest.NewRequest("GET", "/api/promociones/1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en GET /api/promociones/1, obtenido %d", w.Code)
	}

	// 4. Actualizar promoción (PUT)
	newPromo.Value = 60.0
	body, _ = json.Marshal(newPromo)
	req = httptest.NewRequest("PUT", "/api/promociones/1", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en PUT /api/promociones/1, obtenido %d", w.Code)
	}

	// 5. Aplicar decoradores (POST /api/promociones/aplicar)
	applyReq := ApplyPromotionRequest{
		ProductID:    1,
		PromotionIDs: []int64{1, 2}, // Descuento Bienvenida 10% + Bono 25
		CouponCode:   "VIP50",
	}
	body, _ = json.Marshal(applyReq)
	req = httptest.NewRequest("POST", "/api/promociones/aplicar", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en POST /api/promociones/aplicar, obtenido %d", w.Code)
	}

	var res ApplyPromotionResponse
	_ = json.NewDecoder(w.Body).Decode(&res)
	if res.FinalPrice >= res.BasePrice || len(res.Steps) < 2 {
		t.Fatalf("La respuesta del decorador no aplicó descuentos: %+v", res)
	}

	// 6. Eliminar promoción (DELETE)
	req = httptest.NewRequest("DELETE", "/api/promociones/1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en DELETE /api/promociones/1, obtenido %d", w.Code)
	}
}
