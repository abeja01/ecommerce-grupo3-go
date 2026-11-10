package promotions

import (
	"testing"
)

func TestDecoratorPattern(t *testing.T) {
	// 1. Componente base: Producto de 200.00
	base := NewBasePrice("Monitor 27", 200.0)

	if base.GetFinalPrice() != 200.0 {
		t.Fatalf("Esperado 200.0, obtenido %f", base.GetFinalPrice())
	}

	// 2. Decorar con 10% de descuento porcentual -> 200 * 0.9 = 180.0
	decorador1 := NewPercentageDecorator(base, "Cyber Monday 10%", 10.0)
	if decorador1.GetFinalPrice() != 180.0 {
		t.Fatalf("Esperado 180.0, obtenido %f", decorador1.GetFinalPrice())
	}

	// 3. Decorar adicionalmente con Cupón fijo de -S/ 30.0 -> 180 - 30 = 150.0
	decorador2 := NewFixedAmountDecorator(decorador1, "Bono Fidelidad", 30.0)
	if decorador2.GetFinalPrice() != 150.0 {
		t.Fatalf("Esperado 150.0, obtenido %f", decorador2.GetFinalPrice())
	}

	// 4. Decorar con Cupón condicionado con monto mínimo cumplido (mínimo 100)
	decorador3 := NewCouponDecorator(decorador2, "VIP2026", "VIP", "monto_fijo", 20.0, 100.0)
	if decorador3.GetFinalPrice() != 130.0 {
		t.Fatalf("Esperado 130.0, obtenido %f", decorador3.GetFinalPrice())
	}

	// 5. Verificar que el desglose de auditoría tenga 4 pasos
	breakdown := decorador3.GetBreakdown()
	if len(breakdown) != 4 {
		t.Fatalf("Esperado 4 pasos de auditoría en el desglose, obtenido %d", len(breakdown))
	}
}
