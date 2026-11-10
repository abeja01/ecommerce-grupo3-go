package products

import "time"

// Product representa el modelo de datos para un producto (simple o compuesto/combo).
type Product struct {
	ID              int64              `json:"id"`
	CategoryID      *int64             `json:"category_id,omitempty"`
	Name            string             `json:"nombre"`
	Description     string             `json:"descripcion"`
	Price           float64            `json:"precio"`
	Stock           int                `json:"stock"`
	ImageURL        string             `json:"imagen_url"`
	Type            string             `json:"tipo"` // "simple" o "combo"
	DiscountPercent float64            `json:"descuento_porcentaje"`
	Active          bool               `json:"activo"`
	CreatedAt       time.Time          `json:"created_at"`
	Items           []ComboItemPayload `json:"items,omitempty"`
}

// ComboItemPayload define un elemento componente dentro de un combo
type ComboItemPayload struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"nombre,omitempty"`
	Price       float64 `json:"precio_unitario,omitempty"`
	Quantity    int     `json:"cantidad"`
	Subtotal    float64 `json:"subtotal,omitempty"`
}

// ProductInput es el DTO recibido en POST y PUT
type ProductInput struct {
	CategoryID      *int64             `json:"category_id"`
	Name            string             `json:"nombre"`
	Description     string             `json:"descripcion"`
	Price           float64            `json:"precio"`
	Stock           int                `json:"stock"`
	ImageURL        string             `json:"imagen_url"`
	Type            string             `json:"tipo"` // "simple" o "combo"
	DiscountPercent float64            `json:"descuento_porcentaje"`
	Active          bool               `json:"activo"`
	Items           []ComboItemPayload `json:"items"`
}

// CompositeDetailResponse expone la estructura del patrón Composite para la vista/API
type CompositeDetailResponse struct {
	ID              int64              `json:"id"`
	Nombre          string             `json:"nombre"`
	Tipo            string             `json:"tipo"`
	PrecioFinal     float64            `json:"precio_final"`
	PrecioBruto     float64            `json:"precio_bruto"`
	Ahorro          float64            `json:"ahorro"`
	DescuentoPct    float64            `json:"descuento_porcentaje"`
	StockDisponible int                `json:"stock_disponible"`
	Hijos           []ComboItemPayload `json:"elementos_hijos"`
}
