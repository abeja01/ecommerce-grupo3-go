package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce-grupo3/backend/internal/products"
	"ecommerce-grupo3/backend/internal/promotions"
)

func TestFullMonolithAPI(t *testing.T) {
	productRepo := products.NewRepository(nil)
	productHandler := products.NewHandler(productRepo)

	promotionRepo := promotions.NewRepository(nil)
	promotionHandler := promotions.NewHandler(promotionRepo, productRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "proyecto": "E-commerce Grupo 3"})
	})
	productHandler.RegisterRoutes(mux)
	promotionHandler.RegisterRoutes(mux)

	handler := cors(mux)

	// 1. Health check con httptest.NewRecorder (no requiere sockets de red)
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Health check falló: status %d", w.Code)
	}

	// 2. Listar productos (Patrón Composite)
	req = httptest.NewRequest("GET", "/api/productos", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Listar productos falló con status %d", w.Code)
	}
	var prods []products.Product
	_ = json.NewDecoder(w.Body).Decode(&prods)
	if len(prods) == 0 {
		t.Fatalf("Se esperaban productos precargados")
	}

	// 3. Inspeccionar Composite Combo (Pack Gamer ID 6)
	req = httptest.NewRequest("GET", "/api/productos/6/composite", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Composite check falló con status %d", w.Code)
	}
	var comp products.CompositeDetailResponse
	_ = json.NewDecoder(w.Body).Decode(&comp)
	if comp.PrecioFinal <= 0 || len(comp.Hijos) == 0 {
		t.Fatalf("Estructura Composite inválida: %+v", comp)
	}

	// 4. Listar Promociones (Patrón Decorator)
	req = httptest.NewRequest("GET", "/api/promociones", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Listar promociones falló con status %d", w.Code)
	}
	var promos []promotions.Promotion
	_ = json.NewDecoder(w.Body).Decode(&promos)
	if len(promos) == 0 {
		t.Fatalf("Se esperaban promociones precargadas")
	}

	// 5. Aplicar Decoradores encadenados (POST /api/promociones/aplicar)
	applyPayload := map[string]any{
		"producto_id":   6,
		"promocion_ids": []int64{1, 2},
		"codigo_cupon":  "VIP50",
	}
	bodyBytes, _ := json.Marshal(applyPayload)
	req = httptest.NewRequest("POST", "/api/promociones/aplicar", bytes.NewReader(bodyBytes))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Aplicar decoradores falló: status %d (body: %s)", w.Code, w.Body.String())
	}
	var applyRes promotions.ApplyPromotionResponse
	_ = json.NewDecoder(w.Body).Decode(&applyRes)

	if applyRes.FinalPrice >= applyRes.BasePrice {
		t.Fatalf("El precio final con decoradores (%.2f) debería ser menor que el base (%.2f)", applyRes.FinalPrice, applyRes.BasePrice)
	}
	if len(applyRes.Steps) < 2 {
		t.Fatalf("Se esperaban pasos de auditoría de los decoradores, se obtuvieron: %d", len(applyRes.Steps))
	}
}
