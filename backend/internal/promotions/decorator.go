package promotions

import "fmt"

// PriceComponent es la interfaz base del Patrón Decorator.
// Define el contrato para calcular precios y registrar la auditoría paso a paso.
type PriceComponent interface {
	GetDescription() string
	GetBasePrice() float64
	GetFinalPrice() float64
	GetBreakdown() []PriceStep
}

// PriceStep detalla el impacto de cada decorador aplicado
type PriceStep struct {
	StepName        string  `json:"capa_decorador"`
	Description     string  `json:"descripcion"`
	PriceBefore     float64 `json:"precio_anterior"`
	DiscountApplied float64 `json:"descuento_aplicado"`
	PriceAfter      float64 `json:"precio_resultante"`
}

// BasePriceComponent es el componente concreto base (precio original del producto o carrito)
type BasePriceComponent struct {
	Name      string
	BasePrice float64
}

func NewBasePrice(name string, price float64) *BasePriceComponent {
	return &BasePriceComponent{Name: name, BasePrice: price}
}

func (b *BasePriceComponent) GetDescription() string { return b.Name }
func (b *BasePriceComponent) GetBasePrice() float64  { return b.BasePrice }
func (b *BasePriceComponent) GetFinalPrice() float64 { return b.BasePrice }
func (b *BasePriceComponent) GetBreakdown() []PriceStep {
	return []PriceStep{
		{
			StepName:        "Componente Base",
			Description:     fmt.Sprintf("Precio original de '%s'", b.Name),
			PriceBefore:     b.BasePrice,
			DiscountApplied: 0,
			PriceAfter:      b.BasePrice,
		},
	}
}

// BaseDecorator es la clase abstracta / envoltorio base de decoradores
type BaseDecorator struct {
	Wrapped PriceComponent
}

// PercentageDiscountDecorator aplica un descuento porcentual sobre el precio actual
type PercentageDiscountDecorator struct {
	BaseDecorator
	Name       string
	Percentage float64
}

func NewPercentageDecorator(wrapped PriceComponent, name string, pct float64) *PercentageDiscountDecorator {
	return &PercentageDiscountDecorator{
		BaseDecorator: BaseDecorator{Wrapped: wrapped},
		Name:          name,
		Percentage:    pct,
	}
}

func (d *PercentageDiscountDecorator) GetDescription() string {
	return fmt.Sprintf("%s + Decorador Porcentaje (%.1f%%)", d.Wrapped.GetDescription(), d.Percentage)
}

func (d *PercentageDiscountDecorator) GetBasePrice() float64 {
	return d.Wrapped.GetBasePrice()
}

func (d *PercentageDiscountDecorator) GetFinalPrice() float64 {
	prevPrice := d.Wrapped.GetFinalPrice()
	discount := prevPrice * (d.Percentage / 100.0)
	final := prevPrice - discount
	if final < 0 {
		return 0
	}
	return final
}

func (d *PercentageDiscountDecorator) GetBreakdown() []PriceStep {
	steps := d.Wrapped.GetBreakdown()
	prevPrice := d.Wrapped.GetFinalPrice()
	discount := prevPrice * (d.Percentage / 100.0)
	result := prevPrice - discount
	if result < 0 {
		result = 0
	}

	steps = append(steps, PriceStep{
		StepName:        "Decorador Porcentaje",
		Description:     fmt.Sprintf("%s: -%.1f%% (-S/ %.2f)", d.Name, d.Percentage, discount),
		PriceBefore:     prevPrice,
		DiscountApplied: discount,
		PriceAfter:      result,
	})
	return steps
}

// FixedAmountDiscountDecorator aplica un cupón o bono de monto fijo
type FixedAmountDiscountDecorator struct {
	BaseDecorator
	Name   string
	Amount float64
}

func NewFixedAmountDecorator(wrapped PriceComponent, name string, amount float64) *FixedAmountDiscountDecorator {
	return &FixedAmountDiscountDecorator{
		BaseDecorator: BaseDecorator{Wrapped: wrapped},
		Name:          name,
		Amount:        amount,
	}
}

func (d *FixedAmountDiscountDecorator) GetDescription() string {
	return fmt.Sprintf("%s + Decorador Monto Fijo (-S/ %.2f)", d.Wrapped.GetDescription(), d.Amount)
}

func (d *FixedAmountDiscountDecorator) GetBasePrice() float64 {
	return d.Wrapped.GetBasePrice()
}

func (d *FixedAmountDiscountDecorator) GetFinalPrice() float64 {
	prevPrice := d.Wrapped.GetFinalPrice()
	final := prevPrice - d.Amount
	if final < 0 {
		return 0
	}
	return final
}

func (d *FixedAmountDiscountDecorator) GetBreakdown() []PriceStep {
	steps := d.Wrapped.GetBreakdown()
	prevPrice := d.Wrapped.GetFinalPrice()
	discount := d.Amount
	if discount > prevPrice {
		discount = prevPrice
	}
	result := prevPrice - discount

	steps = append(steps, PriceStep{
		StepName:        "Decorador Monto Fijo",
		Description:     fmt.Sprintf("%s: Cupón directo de -S/ %.2f", d.Name, d.Amount),
		PriceBefore:     prevPrice,
		DiscountApplied: discount,
		PriceAfter:      result,
	})
	return steps
}

// CouponDiscountDecorator decora validando un código promocional específico
type CouponDiscountDecorator struct {
	BaseDecorator
	Code      string
	Name      string
	Type      string
	Value     float64
	MinAmount float64
}

func NewCouponDecorator(wrapped PriceComponent, code, name, tipo string, val, minAmount float64) *CouponDiscountDecorator {
	return &CouponDiscountDecorator{
		BaseDecorator: BaseDecorator{Wrapped: wrapped},
		Code:          code,
		Name:          name,
		Type:          tipo,
		Value:         val,
		MinAmount:     minAmount,
	}
}

func (d *CouponDiscountDecorator) GetDescription() string {
	return fmt.Sprintf("%s + Cupón [%s]", d.Wrapped.GetDescription(), d.Code)
}

func (d *CouponDiscountDecorator) GetBasePrice() float64 {
	return d.Wrapped.GetBasePrice()
}

func (d *CouponDiscountDecorator) GetFinalPrice() float64 {
	prev := d.Wrapped.GetFinalPrice()
	if prev < d.MinAmount {
		return prev // No califica por monto mínimo
	}
	var final float64
	if d.Type == "porcentaje" {
		final = prev * (1.0 - (d.Value / 100.0))
	} else {
		final = prev - d.Value
	}
	if final < 0 {
		return 0
	}
	return final
}

func (d *CouponDiscountDecorator) GetBreakdown() []PriceStep {
	steps := d.Wrapped.GetBreakdown()
	prev := d.Wrapped.GetFinalPrice()
	if prev < d.MinAmount {
		steps = append(steps, PriceStep{
			StepName:        "Decorador Cupón (No Aplicado)",
			Description:     fmt.Sprintf("Cupón '%s' requiere monto mínimo de S/ %.2f (monto actual: S/ %.2f)", d.Code, d.MinAmount, prev),
			PriceBefore:     prev,
			DiscountApplied: 0,
			PriceAfter:      prev,
		})
		return steps
	}

	var discount float64
	if d.Type == "porcentaje" {
		discount = prev * (d.Value / 100.0)
	} else {
		discount = d.Value
		if discount > prev {
			discount = prev
		}
	}
	result := prev - discount

	steps = append(steps, PriceStep{
		StepName:        "Decorador Cupón Promocional",
		Description:     fmt.Sprintf("Cupón '%s' (%s): Descuento de S/ %.2f aplicado", d.Code, d.Name, discount),
		PriceBefore:     prev,
		DiscountApplied: discount,
		PriceAfter:      result,
	})
	return steps
}
