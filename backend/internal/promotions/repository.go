package promotions

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db       *pgxpool.Pool
	mu       sync.RWMutex
	memStore map[int64]Promotion
	nextID   int64
}

func NewRepository(db *pgxpool.Pool) *Repository {
	r := &Repository{
		db:       db,
		memStore: make(map[int64]Promotion),
		nextID:   5,
	}
	if db == nil {
		r.seedMemoryStore()
	}
	return r
}

func (r *Repository) List(ctx context.Context) ([]Promotion, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		list := make([]Promotion, 0, len(r.memStore))
		for _, p := range r.memStore {
			list = append(list, p)
		}
		return list, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, codigo, nombre, descripcion, tipo, valor, monto_minimo, activa, created_at, updated_at
		FROM promociones ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Promotion, 0)
	for rows.Next() {
		var p Promotion
		if err := rows.Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.Type, &p.Value, &p.MinAmount, &p.Active, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (Promotion, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		p, exists := r.memStore[id]
		if !exists {
			return Promotion{}, pgx.ErrNoRows
		}
		return p, nil
	}

	var p Promotion
	err := r.db.QueryRow(ctx, `
		SELECT id, codigo, nombre, descripcion, tipo, valor, monto_minimo, activa, created_at, updated_at
		FROM promociones WHERE id = $1`, id).
		Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.Type, &p.Value, &p.MinAmount, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) GetByCode(ctx context.Context, code string) (Promotion, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for _, p := range r.memStore {
			if strings.EqualFold(p.Code, code) {
				return p, nil
			}
		}
		return Promotion{}, pgx.ErrNoRows
	}

	var p Promotion
	err := r.db.QueryRow(ctx, `
		SELECT id, codigo, nombre, descripcion, tipo, valor, monto_minimo, activa, created_at, updated_at
		FROM promociones WHERE UPPER(codigo) = $1`, code).
		Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.Type, &p.Value, &p.MinAmount, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) Create(ctx context.Context, in PromotionInput) (Promotion, error) {
	in.Code = strings.ToUpper(strings.TrimSpace(in.Code))
	now := time.Now()

	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		for _, ex := range r.memStore {
			if strings.EqualFold(ex.Code, in.Code) {
				return Promotion{}, fmt.Errorf("ya existe una promoción con el código '%s'", in.Code)
			}
		}
		p := Promotion{
			ID:          r.nextID,
			Code:        in.Code,
			Name:        in.Name,
			Description: in.Description,
			Type:        in.Type,
			Value:       in.Value,
			MinAmount:   in.MinAmount,
			Active:      in.Active,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		r.nextID++
		r.memStore[p.ID] = p
		return p, nil
	}

	var p Promotion
	err := r.db.QueryRow(ctx, `
		INSERT INTO promociones(codigo, nombre, descripcion, tipo, valor, monto_minimo, activa)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, codigo, nombre, descripcion, tipo, valor, monto_minimo, activa, created_at, updated_at`,
		in.Code, in.Name, in.Description, in.Type, in.Value, in.MinAmount, in.Active).
		Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.Type, &p.Value, &p.MinAmount, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *Repository) Update(ctx context.Context, id int64, in PromotionInput) (Promotion, error) {
	in.Code = strings.ToUpper(strings.TrimSpace(in.Code))
	now := time.Now()

	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		p, exists := r.memStore[id]
		if !exists {
			return Promotion{}, pgx.ErrNoRows
		}
		p.Code = in.Code
		p.Name = in.Name
		p.Description = in.Description
		p.Type = in.Type
		p.Value = in.Value
		p.MinAmount = in.MinAmount
		p.Active = in.Active
		p.UpdatedAt = now
		r.memStore[id] = p
		return p, nil
	}

	var p Promotion
	err := r.db.QueryRow(ctx, `
		UPDATE promociones
		SET codigo=$1, nombre=$2, descripcion=$3, tipo=$4, valor=$5, monto_minimo=$6, activa=$7, updated_at=NOW()
		WHERE id=$8
		RETURNING id, codigo, nombre, descripcion, tipo, valor, monto_minimo, activa, created_at, updated_at`,
		in.Code, in.Name, in.Description, in.Type, in.Value, in.MinAmount, in.Active, id).
		Scan(&p.ID, &p.Code, &p.Name, &p.Description, &p.Type, &p.Value, &p.MinAmount, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	return p, err
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

	result, err := r.db.Exec(ctx, `DELETE FROM promociones WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) seedMemoryStore() {
	now := time.Now()
	r.memStore[1] = Promotion{ID: 1, Code: "BIENVENIDA10", Name: "Descuento Bienvenida 10%", Description: "10% de descuento en cualquier compra", Type: "porcentaje", Value: 10.0, MinAmount: 0.0, Active: true, CreatedAt: now, UpdatedAt: now}
	r.memStore[2] = Promotion{ID: 2, Code: "BONO25", Name: "Cupón Fijo Grupo 3 - S/ 25", Description: "Descuento de S/ 25 en compras > S/ 100", Type: "monto_fijo", Value: 25.0, MinAmount: 100.0, Active: true, CreatedAt: now, UpdatedAt: now}
	r.memStore[3] = Promotion{ID: 3, Code: "CYBER2026", Name: "Cyber Promo 20%", Description: "20% en tecnología seleccionada", Type: "porcentaje", Value: 20.0, MinAmount: 150.0, Active: true, CreatedAt: now, UpdatedAt: now}
	r.memStore[4] = Promotion{ID: 4, Code: "VIP50", Name: "Voucher VIP - S/ 50", Description: "Cupón de S/ 50 para compras > S/ 500", Type: "cupon", Value: 50.0, MinAmount: 500.0, Active: true, CreatedAt: now, UpdatedAt: now}
}

func ValidateInput(in PromotionInput) error {
	if strings.TrimSpace(in.Code) == "" {
		return fmt.Errorf("el código promocional es obligatorio")
	}
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("el nombre de la promoción es obligatorio")
	}
	if in.Type != "porcentaje" && in.Type != "monto_fijo" && in.Type != "cupon" {
		return fmt.Errorf("el tipo de promoción debe ser 'porcentaje', 'monto_fijo' o 'cupon'")
	}
	if in.Value < 0 {
		return fmt.Errorf("el valor del descuento no puede ser negativo")
	}
	if in.Type == "porcentaje" && in.Value > 100 {
		return fmt.Errorf("el porcentaje de descuento no puede ser superior al 100%%")
	}
	if in.MinAmount < 0 {
		return fmt.Errorf("el monto mínimo no puede ser negativo")
	}
	return nil
}
