package products

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db       *pgxpool.Pool
	mu       sync.RWMutex
	memStore map[int64]Product
	nextID   int64
}

func NewRepository(db *pgxpool.Pool) *Repository {
	r := &Repository{
		db:       db,
		memStore: make(map[int64]Product),
		nextID:   8,
	}
	if db == nil {
		r.seedMemoryStore()
	}
	return r
}

func (r *Repository) List(ctx context.Context) ([]Product, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		list := make([]Product, 0, len(r.memStore))
		for _, p := range r.memStore {
			list = append(list, r.enrichWithComposite(p))
		}
		return list, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje, activo, created_at 
		FROM productos ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Product, 0)
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.Type, &p.DiscountPercent, &p.Active, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range items {
		if items[i].Type == "combo" {
			items[i] = r.loadComboDetailsDB(ctx, items[i])
		}
	}
	return items, nil
}

func (r *Repository) Get(ctx context.Context, id int64) (Product, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		p, exists := r.memStore[id]
		if !exists {
			return Product{}, pgx.ErrNoRows
		}
		return r.enrichWithComposite(p), nil
	}

	var p Product
	err := r.db.QueryRow(ctx, `
		SELECT id, category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje, activo, created_at 
		FROM productos WHERE id = $1`, id).
		Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.Type, &p.DiscountPercent, &p.Active, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	if p.Type == "combo" {
		p = r.loadComboDetailsDB(ctx, p)
	}
	return p, nil
}

func (r *Repository) Create(ctx context.Context, in ProductInput) (Product, error) {
	if in.Type == "" {
		in.Type = "simple"
	}

	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		p := Product{
			ID:              r.nextID,
			CategoryID:      in.CategoryID,
			Name:            in.Name,
			Description:     in.Description,
			Price:           in.Price,
			Stock:           in.Stock,
			ImageURL:        in.ImageURL,
			Type:            in.Type,
			DiscountPercent: in.DiscountPercent,
			Active:          in.Active,
			CreatedAt:       time.Now(),
			Items:           in.Items,
		}
		r.nextID++
		r.memStore[p.ID] = p
		return r.enrichWithComposite(p), nil
	}

	var p Product
	err := r.db.QueryRow(ctx, `
		INSERT INTO productos(category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje, activo)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje, activo, created_at`,
		in.CategoryID, in.Name, in.Description, in.Price, in.Stock, in.ImageURL, in.Type, in.DiscountPercent, in.Active).
		Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.Type, &p.DiscountPercent, &p.Active, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	if in.Type == "combo" && len(in.Items) > 0 {
		for _, item := range in.Items {
			_, _ = r.db.Exec(ctx, `
				INSERT INTO combo_items(combo_id, item_id, cantidad) 
				VALUES($1, $2, $3) 
				ON CONFLICT (combo_id, item_id) DO UPDATE SET cantidad = EXCLUDED.cantidad`,
				p.ID, item.ProductID, item.Quantity)
		}
		p = r.loadComboDetailsDB(ctx, p)
	}

	return p, nil
}

func (r *Repository) Update(ctx context.Context, id int64, in ProductInput) (Product, error) {
	if in.Type == "" {
		in.Type = "simple"
	}

	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		p, exists := r.memStore[id]
		if !exists {
			return Product{}, pgx.ErrNoRows
		}
		p.CategoryID = in.CategoryID
		p.Name = in.Name
		p.Description = in.Description
		p.Price = in.Price
		p.Stock = in.Stock
		p.ImageURL = in.ImageURL
		p.Type = in.Type
		p.DiscountPercent = in.DiscountPercent
		p.Active = in.Active
		p.Items = in.Items
		r.memStore[id] = p
		return r.enrichWithComposite(p), nil
	}

	var p Product
	err := r.db.QueryRow(ctx, `
		UPDATE productos 
		SET category_id=$1, nombre=$2, descripcion=$3, precio=$4, stock=$5, imagen_url=$6, tipo=$7, descuento_porcentaje=$8, activo=$9, updated_at=NOW()
		WHERE id=$10
		RETURNING id, category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje, activo, created_at`,
		in.CategoryID, in.Name, in.Description, in.Price, in.Stock, in.ImageURL, in.Type, in.DiscountPercent, in.Active, id).
		Scan(&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.Type, &p.DiscountPercent, &p.Active, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	if in.Type == "combo" {
		_, _ = r.db.Exec(ctx, `DELETE FROM combo_items WHERE combo_id = $1`, id)
		for _, item := range in.Items {
			_, _ = r.db.Exec(ctx, `
				INSERT INTO combo_items(combo_id, item_id, cantidad) 
				VALUES($1, $2, $3)`,
				id, item.ProductID, item.Quantity)
		}
		p = r.loadComboDetailsDB(ctx, p)
	}

	return p, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, exists := r.memStore[id]; !exists {
			return pgx.ErrNoRows
		}
		delete(r.memStore, id)
		return nil
	}

	result, err := r.db.Exec(ctx, `DELETE FROM productos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) loadComboDetailsDB(ctx context.Context, combo Product) Product {
	rows, err := r.db.Query(ctx, `
		SELECT ci.item_id, p.nombre, p.precio, ci.cantidad, (p.precio * ci.cantidad) AS subtotal, p.stock, p.descripcion
		FROM combo_items ci
		JOIN productos p ON p.id = ci.item_id
		WHERE ci.combo_id = $1`, combo.ID)
	if err != nil {
		return combo
	}
	defer rows.Close()

	var childDetails []ChildDetail
	var items []ComboItemPayload

	for rows.Next() {
		var itemID int64
		var name, desc string
		var unitPrice float64
		var qty, stock int
		var subtotal float64
		if err := rows.Scan(&itemID, &name, &unitPrice, &qty, &subtotal, &stock, &desc); err == nil {
			leaf := LeafProduct{ID: itemID, Name: name, Description: desc, Price: unitPrice, Stock: stock}
			childDetails = append(childDetails, ChildDetail{
				Component:   leaf,
				Cantidad:    qty,
				PrecioTotal: subtotal,
			})
			items = append(items, ComboItemPayload{
				ProductID:   itemID,
				ProductName: name,
				Price:       unitPrice,
				Quantity:    qty,
				Subtotal:    subtotal,
			})
		}
	}

	composite := CompositeProduct{
		ID:              combo.ID,
		Name:            combo.Name,
		Description:     combo.Description,
		DiscountPercent: combo.DiscountPercent,
		Children:        childDetails,
	}

	combo.Price = composite.GetPrice()
	combo.Stock = composite.GetStock()
	combo.Items = items
	return combo
}

func (r *Repository) GetPriceAndName(id int64) (string, float64, error) {
	p, err := r.Get(context.Background(), id)
	if err != nil {
		return "", 0, err
	}
	return p.Name, p.Price, nil
}

func (r *Repository) enrichWithComposite(p Product) Product {
	if p.Type != "combo" {
		return p
	}
	var childDetails []ChildDetail
	var items []ComboItemPayload
	for _, it := range p.Items {
		child, exists := r.memStore[it.ProductID]
		if !exists {
			continue
		}
		leaf := LeafProduct{ID: child.ID, Name: child.Name, Description: child.Description, Price: child.Price, Stock: child.Stock}
		subtotal := child.Price * float64(it.Quantity)
		childDetails = append(childDetails, ChildDetail{
			Component:   leaf,
			Cantidad:    it.Quantity,
			PrecioTotal: subtotal,
		})
		items = append(items, ComboItemPayload{
			ProductID:   child.ID,
			ProductName: child.Name,
			Price:       child.Price,
			Quantity:    it.Quantity,
			Subtotal:    subtotal,
		})
	}

	comp := CompositeProduct{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		DiscountPercent: p.DiscountPercent,
		Children:        childDetails,
	}
	p.Price = comp.GetPrice()
	p.Stock = comp.GetStock()
	p.Items = items
	return p
}

func (r *Repository) seedMemoryStore() {
	cat1, cat2, cat3, cat4 := int64(1), int64(2), int64(3), int64(4)
	r.memStore[1] = Product{ID: 1, CategoryID: &cat1, Name: "Laptop Pro 14", Description: "Laptop para trabajo y código", Price: 3299.90, Stock: 10, ImageURL: "https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=900", Type: "simple", Active: true, CreatedAt: time.Now()}
	r.memStore[2] = Product{ID: 2, CategoryID: &cat2, Name: "Mouse Inalámbrico RGB", Description: "Mouse ergonómico de precisión", Price: 79.90, Stock: 30, ImageURL: "https://images.unsplash.com/photo-1527814050087-3793815479db?w=900", Type: "simple", Active: true, CreatedAt: time.Now()}
	r.memStore[3] = Product{ID: 3, CategoryID: &cat2, Name: "Teclado Mecánico RGB", Description: "Teclado mecánico retroiluminado switch azul", Price: 189.90, Stock: 20, ImageURL: "https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=900", Type: "simple", Active: true, CreatedAt: time.Now()}
	r.memStore[4] = Product{ID: 4, CategoryID: &cat3, Name: "Audífonos Bluetooth Pro", Description: "Audífonos con cancelación de ruido", Price: 149.90, Stock: 25, ImageURL: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=900", Type: "simple", Active: true, CreatedAt: time.Now()}
	r.memStore[5] = Product{ID: 5, CategoryID: &cat2, Name: "Monitor IPS 27\" 4K", Description: "Monitor ultra nítido 4K", Price: 899.90, Stock: 12, ImageURL: "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=900", Type: "simple", Active: true, CreatedAt: time.Now()}

	r.memStore[6] = Product{
		ID:              6,
		CategoryID:      &cat4,
		Name:            "Pack Periféricos Gamer Pro",
		Description:     "Combo completo: Teclado Mecánico + Mouse RGB + Audífonos Bluetooth con 10% de descuento",
		Type:            "combo",
		DiscountPercent: 10.0,
		Active:          true,
		ImageURL:        "https://images.unsplash.com/photo-1542751371-adc38448a05e?w=900",
		CreatedAt:       time.Now(),
		Items: []ComboItemPayload{
			{ProductID: 3, Quantity: 1},
			{ProductID: 2, Quantity: 1},
			{ProductID: 4, Quantity: 1},
		},
	}

	r.memStore[7] = Product{
		ID:              7,
		CategoryID:      &cat4,
		Name:            "Combo Workstation Ultimate",
		Description:     "Setup integral: Laptop Pro 14 + Monitor 27\" 4K + Teclado + Mouse con 15% de descuento",
		Type:            "combo",
		DiscountPercent: 15.0,
		Active:          true,
		ImageURL:        "https://images.unsplash.com/photo-1518770660439-4636190af475?w=900",
		CreatedAt:       time.Now(),
		Items: []ComboItemPayload{
			{ProductID: 1, Quantity: 1},
			{ProductID: 5, Quantity: 1},
			{ProductID: 3, Quantity: 1},
			{ProductID: 2, Quantity: 1},
		},
	}
}

func ValidateInput(in ProductInput) error {
	if in.Name == "" {
		return fmt.Errorf("el nombre es obligatorio")
	}
	if in.Type == "combo" {
		if in.DiscountPercent < 0 || in.DiscountPercent > 100 {
			return fmt.Errorf("el descuento de combo debe estar entre 0%% y 100%%")
		}
		if len(in.Items) == 0 {
			return fmt.Errorf("un combo debe incluir al menos un producto componente")
		}
	} else {
		if in.Price < 0 {
			return fmt.Errorf("el precio no puede ser negativo")
		}
		if in.Stock < 0 {
			return fmt.Errorf("el stock no puede ser negativo")
		}
	}
	return nil
}
