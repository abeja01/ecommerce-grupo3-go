package promotions

import "time"

// Promotion representa la entidad de Promoción en la base de datos
type Promotion struct {
	ID          int64     `json:"id"`
	Code        string    `json:"codigo"`
	Name        string    `json:"nombre"`
	Description string    `json:"descripcion"`
	Type        string    `json:"tipo"` // "porcentaje", "monto_fijo", "cupon"
	Value       float64   `json:"valor"`
	MinAmount   float64   `json:"monto_minimo"`
	Active      bool      `json:"activa"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PromotionInput es el DTO para crear o editar una promoción
type PromotionInput struct {
	Code        string  `json:"codigo"`
	Name        string  `json:"nombre"`
	Description string  `json:"descripcion"`
	Type        string  `json:"tipo"`
	Value       float64 `json:"valor"`
	MinAmount   float64 `json:"monto_minimo"`
	Active      bool    `json:"activa"`
}

// ApplyPromotionRequest es la carga para probar el encadenamiento de decoradores
type ApplyPromotionRequest struct {
	ProductID    int64   `json:"producto_id"`
	CustomAmount float64 `json:"monto_personalizado"` // Si se evalúa monto directo de carrito
	PromotionIDs []int64 `json:"promocion_ids"`
	CouponCode   string  `json:"codigo_cupon"`
}

// ApplyPromotionResponse expone el cálculo final y cada capa de decorador aplicada
type ApplyPromotionResponse struct {
	ProductID       int64       `json:"producto_id"`
	ProductName     string      `json:"nombre_producto"`
	BasePrice       float64     `json:"precio_base"`
	TotalDiscount   float64     `json:"descuento_total"`
	FinalPrice      float64     `json:"precio_final"`
	DecoratorsCount int         `json:"total_decoradores_aplicados"`
	Steps           []PriceStep `json:"pasos_decorador"`
}
