package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL no está configurada")
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL inválida: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("no se pudo conectar con el servidor PostgreSQL: %w", err)
	}

	log.Println("[INFO] Conexión exitosa a PostgreSQL.")
	if err := AutoMigrate(ctx, pool); err != nil {
		log.Printf("[WARN] Error al verificar/migrar tablas: %v", err)
	}

	return pool, nil
}

// AutoMigrate asegura que las tablas requeridas por los dos CRUDs existan
func AutoMigrate(ctx context.Context, db *pgxpool.Pool) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS categorias (
		id BIGSERIAL PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL UNIQUE,
		descripcion TEXT,
		activa BOOLEAN NOT NULL DEFAULT TRUE
	);

	CREATE TABLE IF NOT EXISTS productos (
		id BIGSERIAL PRIMARY KEY,
		category_id BIGINT REFERENCES categorias(id) ON DELETE SET NULL,
		nombre VARCHAR(160) NOT NULL,
		descripcion TEXT NOT NULL DEFAULT '',
		precio NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (precio >= 0),
		stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
		imagen_url TEXT NOT NULL DEFAULT '',
		tipo VARCHAR(20) NOT NULL DEFAULT 'simple' CHECK (tipo IN ('simple', 'combo')),
		descuento_porcentaje NUMERIC(5,2) NOT NULL DEFAULT 0,
		activo BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS combo_items (
		combo_id BIGINT NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
		item_id BIGINT NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
		cantidad INTEGER NOT NULL DEFAULT 1 CHECK (cantidad > 0),
		PRIMARY KEY (combo_id, item_id)
	);

	CREATE TABLE IF NOT EXISTS promociones (
		id BIGSERIAL PRIMARY KEY,
		codigo VARCHAR(50) NOT NULL UNIQUE,
		nombre VARCHAR(120) NOT NULL,
		descripcion TEXT NOT NULL DEFAULT '',
		tipo VARCHAR(30) NOT NULL CHECK (tipo IN ('porcentaje', 'monto_fijo', 'cupon')),
		valor NUMERIC(12,2) NOT NULL DEFAULT 0,
		monto_minimo NUMERIC(12,2) NOT NULL DEFAULT 0,
		activa BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	_, err := db.Exec(ctx, ddl)
	return err
}
