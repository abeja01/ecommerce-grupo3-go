package products

// Component define la interfaz común del Patrón Composite.
// Permite que tanto productos individuales (hojas) como combos/paquetes (compuestos)
// sean tratados de manera uniforme por los clientes y la API.
type Component interface {
	GetID() int64
	GetName() string
	GetPrice() float64
	GetStock() int
	GetDescription() string
	GetType() string
	GetChildren() []ChildDetail
}

// ChildDetail representa un elemento hijo dentro de un compuesto (Composite)
type ChildDetail struct {
	Component   Component `json:"componente"`
	Cantidad    int       `json:"cantidad"`
	PrecioTotal float64   `json:"precio_total"`
}

// LeafProduct representa la hoja del patrón Composite (un producto simple).
type LeafProduct struct {
	ID          int64   `json:"id"`
	Name        string  `json:"nombre"`
	Description string  `json:"descripcion"`
	Price       float64 `json:"precio"`
	Stock       int     `json:"stock"`
}

func (l LeafProduct) GetID() int64                { return l.ID }
func (l LeafProduct) GetName() string             { return l.Name }
func (l LeafProduct) GetPrice() float64           { return l.Price }
func (l LeafProduct) GetStock() int               { return l.Stock }
func (l LeafProduct) GetDescription() string      { return l.Description }
func (l LeafProduct) GetType() string             { return "simple" }
func (l LeafProduct) GetChildren() []ChildDetail { return nil }

// CompositeProduct representa un contenedor compuesto (Combo o Pack de productos).
// Puede contener tanto hojas (productos simples) como otros compuestos.
type CompositeProduct struct {
	ID              int64         `json:"id"`
	Name            string        `json:"nombre"`
	Description     string        `json:"descripcion"`
	DiscountPercent float64       `json:"descuento_porcentaje"`
	Children        []ChildDetail `json:"elementos"`
}

func (c CompositeProduct) GetID() int64           { return c.ID }
func (c CompositeProduct) GetName() string        { return c.Name }
func (c CompositeProduct) GetDescription() string { return c.Description }
func (c CompositeProduct) GetType() string        { return "combo" }
func (c CompositeProduct) GetChildren() []ChildDetail {
	return c.Children
}

// GetPrice calcula recursivamente el precio total de todos los hijos en el Composite
// aplicando el porcentaje de descuento del paquete si corresponde.
func (c CompositeProduct) GetPrice() float64 {
	var rawSum float64
	for _, child := range c.Children {
		rawSum += child.Component.GetPrice() * float64(child.Cantidad)
	}
	if c.DiscountPercent > 0 && c.DiscountPercent <= 100 {
		return rawSum * (1.0 - (c.DiscountPercent / 100.0))
	}
	return rawSum
}

// GetRawPrice retorna la suma bruta sin el descuento del combo
func (c CompositeProduct) GetRawPrice() float64 {
	var rawSum float64
	for _, child := range c.Children {
		rawSum += child.Component.GetPrice() * float64(child.Cantidad)
	}
	return rawSum
}

// GetStock calcula el stock disponible para el combo, que corresponde a la cantidad
// máxima de combos que se pueden ensamblar según el stock de sus componentes.
func (c CompositeProduct) GetStock() int {
	if len(c.Children) == 0 {
		return 0
	}
	minStock := -1
	for _, child := range c.Children {
		childStock := child.Component.GetStock()
		if child.Cantidad <= 0 {
			continue
		}
		possible := childStock / child.Cantidad
		if minStock == -1 || possible < minStock {
			minStock = possible
		}
	}
	if minStock < 0 {
		return 0
	}
	return minStock
}
