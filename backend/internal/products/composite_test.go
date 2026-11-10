package products

import (
	"testing"
)

func TestCompositePattern(t *testing.T) {
	// 1. Crear hojas (LeafProduct)
	teclado := LeafProduct{ID: 1, Name: "Teclado Mecánico", Price: 100.0, Stock: 10}
	mouse := LeafProduct{ID: 2, Name: "Mouse Gamer", Price: 50.0, Stock: 5}

	if teclado.GetPrice() != 100.0 {
		t.Fatalf("Esperado 100.0, obtenido %f", teclado.GetPrice())
	}
	if mouse.GetStock() != 5 {
		t.Fatalf("Esperado 5, obtenido %d", mouse.GetStock())
	}

	// 2. Crear compuesto (Combo) con 10% de descuento
	combo := CompositeProduct{
		ID:              10,
		Name:            "Pack Gamer Básico",
		DiscountPercent: 10.0,
		Children: []ChildDetail{
			{Component: teclado, Cantidad: 1, PrecioTotal: 100.0},
			{Component: mouse, Cantidad: 1, PrecioTotal: 50.0},
		},
	}

	// Precio bruto = 100 + 50 = 150
	// Con 10% descuento = 150 * 0.90 = 135.0
	expectedPrice := 135.0
	if combo.GetPrice() != expectedPrice {
		t.Fatalf("Esperado precio con descuento %f, obtenido %f", expectedPrice, combo.GetPrice())
	}

	// Stock disponible debe ser min(10/1, 5/1) = 5
	expectedStock := 5
	if combo.GetStock() != expectedStock {
		t.Fatalf("Esperado stock %d, obtenido %d", expectedStock, combo.GetStock())
	}

	// 3. Compuesto anidado (Combo dentro de otro Combo)
	audifonos := LeafProduct{ID: 3, Name: "Audífonos", Price: 50.0, Stock: 8}
	comboMaster := CompositeProduct{
		ID:              20,
		Name:            "Super Pack",
		DiscountPercent: 0,
		Children: []ChildDetail{
			{Component: combo, Cantidad: 1, PrecioTotal: 135.0},
			{Component: audifonos, Cantidad: 1, PrecioTotal: 50.0},
		},
	}

	expectedMasterPrice := 135.0 + 50.0 // 185.0
	if comboMaster.GetPrice() != expectedMasterPrice {
		t.Fatalf("Esperado precio compuesto %f, obtenido %f", expectedMasterPrice, comboMaster.GetPrice())
	}
}
