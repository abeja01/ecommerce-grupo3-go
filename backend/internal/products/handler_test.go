package products

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouter() *http.ServeMux {
	repo := NewRepository(nil) // repositorio en memoria
	h := NewHandler(repo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestProductsCRUD(t *testing.T) {
	mux := setupTestRouter()

	// 1. Listar productos (GET)
	req := httptest.NewRequest("GET", "/api/productos", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en GET /api/productos, obtenido %d", w.Code)
	}

	// 2. Crear producto simple (POST)
	newProd := ProductInput{
		Name:        "Webcam 1080p Pro",
		Description: "Cámara web para streaming",
		Price:       120.0,
		Stock:       15,
		Type:        "simple",
		Active:      true,
	}
	body, _ := json.Marshal(newProd)
	req = httptest.NewRequest("POST", "/api/productos", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Esperado 201 en POST /api/productos, obtenido %d: %s", w.Code, w.Body.String())
	}

	var created Product
	_ = json.NewDecoder(w.Body).Decode(&created)
	if created.ID == 0 || created.Name != "Webcam 1080p Pro" {
		t.Fatalf("Producto creado no coincide: %+v", created)
	}

	// 3. Consultar producto por ID (GET)
	req = httptest.NewRequest("GET", "/api/productos/6", nil) // Combo 6
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en GET /api/productos/6, obtenido %d", w.Code)
	}

	// 4. Detalle Composite (GET /api/productos/6/composite)
	req = httptest.NewRequest("GET", "/api/productos/6/composite", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en GET /api/productos/6/composite, obtenido %d", w.Code)
	}

	// 5. Actualizar producto (PUT)
	newProd.Price = 135.0
	body, _ = json.Marshal(newProd)
	req = httptest.NewRequest("PUT", "/api/productos/1", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en PUT /api/productos/1, obtenido %d", w.Code)
	}

	// 6. Eliminar producto (DELETE)
	req = httptest.NewRequest("DELETE", "/api/productos/1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Esperado 200 en DELETE /api/productos/1, obtenido %d", w.Code)
	}
}
