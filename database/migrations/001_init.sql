-- =======================================================
-- Base de Datos: E-commerce Grupo 3
-- Arquitectura Monolítica: Go + PostgreSQL + React
-- CRUD 1: Productos (Patrón Composite - Productos Simples y Combos)
-- CRUD 2: Promociones (Patrón Decorator - Descuentos Acumulables)
-- =======================================================

CREATE TABLE IF NOT EXISTS usuarios (
  id BIGSERIAL PRIMARY KEY,
  nombre VARCHAR(120) NOT NULL,
  email VARCHAR(180) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  rol VARCHAR(20) NOT NULL DEFAULT 'cliente' CHECK (rol IN ('cliente','admin')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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
  descuento_porcentaje NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (descuento_porcentaje >= 0 AND descuento_porcentaje <= 100),
  activo BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabla para el Patrón Composite: Relación de ítems contenidos en un Combo/Pack
CREATE TABLE IF NOT EXISTS combo_items (
  combo_id BIGINT NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
  item_id BIGINT NOT NULL REFERENCES productos(id) ON DELETE CASCADE,
  cantidad INTEGER NOT NULL DEFAULT 1 CHECK (cantidad > 0),
  PRIMARY KEY (combo_id, item_id)
);

-- Tabla para el Patrón Decorator: CRUD de Promociones y Descuentos
CREATE TABLE IF NOT EXISTS promociones (
  id BIGSERIAL PRIMARY KEY,
  codigo VARCHAR(50) NOT NULL UNIQUE,
  nombre VARCHAR(120) NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
  tipo VARCHAR(30) NOT NULL CHECK (tipo IN ('porcentaje', 'monto_fijo', 'cupon')),
  valor NUMERIC(12,2) NOT NULL CHECK (valor >= 0),
  monto_minimo NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (monto_minimo >= 0),
  activa BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pedidos (
  id BIGSERIAL PRIMARY KEY,
  usuario_id BIGINT NOT NULL REFERENCES usuarios(id),
  estado VARCHAR(20) NOT NULL DEFAULT 'pendiente' CHECK (estado IN ('pendiente','pagado','preparando','enviado','entregado','cancelado')),
  subtotal NUMERIC(12,2) NOT NULL CHECK (subtotal >= 0),
  total NUMERIC(12,2) NOT NULL CHECK (total >= 0),
  direccion_entrega TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS detalle_pedidos (
  id BIGSERIAL PRIMARY KEY,
  pedido_id BIGINT NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
  producto_id BIGINT NOT NULL REFERENCES productos(id),
  cantidad INTEGER NOT NULL CHECK (cantidad > 0),
  precio_unitario NUMERIC(12,2) NOT NULL CHECK (precio_unitario >= 0),
  subtotal NUMERIC(12,2) GENERATED ALWAYS AS (cantidad * precio_unitario) STORED
);

CREATE TABLE IF NOT EXISTS pagos (
  id BIGSERIAL PRIMARY KEY,
  pedido_id BIGINT NOT NULL UNIQUE REFERENCES pedidos(id),
  metodo VARCHAR(40) NOT NULL,
  estado VARCHAR(20) NOT NULL DEFAULT 'simulado' CHECK (estado IN ('simulado','aprobado','rechazado')),
  monto NUMERIC(12,2) NOT NULL CHECK (monto >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índices de optimización
CREATE INDEX IF NOT EXISTS idx_productos_categoria ON productos(category_id);
CREATE INDEX IF NOT EXISTS idx_productos_tipo ON productos(tipo);
CREATE INDEX IF NOT EXISTS idx_combo_items_combo ON combo_items(combo_id);
CREATE INDEX IF NOT EXISTS idx_promociones_codigo ON promociones(codigo);
CREATE INDEX IF NOT EXISTS idx_pedidos_usuario ON pedidos(usuario_id);

-- =======================================================
-- DATOS SEMILLA (SEED DATA)
-- =======================================================

-- 1. Categorías
INSERT INTO categorias (id, nombre, descripcion) VALUES
  (1, 'Computadoras', 'Laptops y equipos de escritorio'),
  (2, 'Accesorios', 'Periféricos y accesorios tecnológicos'),
  (3, 'Audio', 'Audífonos y parlantes'),
  (4, 'Combos Especiales', 'Paquetes de productos armados con descuento (Composite)')
ON CONFLICT (id) DO NOTHING;

-- 2. Productos Simples (Hojas del patrón Composite)
INSERT INTO productos (id, category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje) VALUES
  (1, 1, 'Laptop Pro 14', 'Laptop de alto rendimiento para trabajo y programación', 3299.90, 10, 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=900', 'simple', 0),
  (2, 2, 'Mouse Inalámbrico RGB', 'Mouse ergonómico de alta precisión con sensor óptico', 79.90, 30, 'https://images.unsplash.com/photo-1527814050087-3793815479db?w=900', 'simple', 0),
  (3, 2, 'Teclado Mecánico RGB', 'Teclado mecánico retroiluminado switch azul para desarrolladores', 189.90, 20, 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=900', 'simple', 0),
  (4, 3, 'Audífonos Bluetooth Pro', 'Audífonos inalámbricos con cancelación de ruido activa', 149.90, 25, 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=900', 'simple', 0),
  (5, 2, 'Monitor IPS 27" 4K', 'Monitor profesional ultra nítido para diseño y código', 899.90, 12, 'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=900', 'simple', 0)
ON CONFLICT (id) DO NOTHING;

-- 3. Productos Compuestos (Combos / Composite Pattern)
-- Combo 1: Pack Periféricos Pro (Mouse + Teclado + Audífonos) con 10% de descuento bundle
INSERT INTO productos (id, category_id, nombre, descripcion, precio, stock, imagen_url, tipo, descuento_porcentaje) VALUES
  (6, 4, 'Pack Periféricos Gamer Pro', 'Combo completo: Teclado Mecánico + Mouse RGB + Audífonos Bluetooth con 10% de descuento', 0, 0, 'https://images.unsplash.com/photo-1542751371-adc38448a05e?w=900', 'combo', 10.00),
  (7, 4, 'Combo Workstation Ultimate', 'Setup integral: Laptop Pro 14 + Monitor 27" 4K + Teclado + Mouse con 15% de descuento', 0, 0, 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=900', 'combo', 15.00)
ON CONFLICT (id) DO NOTHING;

-- Relaciones de los Combos (Patrón Composite)
-- Pack 6 (Periféricos): 1 Teclado(3) + 1 Mouse(2) + 1 Audífonos(4)
INSERT INTO combo_items (combo_id, item_id, cantidad) VALUES
  (6, 3, 1),
  (6, 2, 1),
  (6, 4, 1)
ON CONFLICT (combo_id, item_id) DO NOTHING;

-- Pack 7 (Workstation): 1 Laptop(1) + 1 Monitor(5) + 1 Teclado(3) + 1 Mouse(2)
INSERT INTO combo_items (combo_id, item_id, cantidad) VALUES
  (7, 1, 1),
  (7, 5, 1),
  (7, 3, 1),
  (7, 2, 1)
ON CONFLICT (combo_id, item_id) DO NOTHING;

-- 4. Promociones (Patrón Decorator)
INSERT INTO promociones (id, codigo, nombre, descripcion, tipo, valor, monto_minimo, activa) VALUES
  (1, 'BIENVENIDA10', 'Descuento de Bienvenida 10%', '10% de descuento en cualquier compra sin monto mínimo', 'porcentaje', 10.00, 0.00, true),
  (2, 'BONO25', 'Cupón Fijo Grupo 3 - S/ 25', 'Descuento directo de S/ 25.00 en compras mayores a S/ 100', 'monto_fijo', 25.00, 100.00, true),
  (3, 'CYBER2026', 'Cyber Promo Especial 20%', '20% de descuento en tecnología seleccionada por temporada', 'porcentaje', 20.00, 150.00, true),
  (4, 'VIP50', 'Voucher VIP Tecnológico - S/ 50', 'Descuento exclusivo de S/ 50.00 para compras mayores a S/ 500', 'cupon', 50.00, 500.00, true)
ON CONFLICT (id) DO NOTHING;

-- Ajustar secuencias de ID
SELECT setval('categorias_id_seq', (SELECT COALESCE(MAX(id), 1) FROM categorias));
SELECT setval('productos_id_seq', (SELECT COALESCE(MAX(id), 1) FROM productos));
SELECT setval('promociones_id_seq', (SELECT COALESCE(MAX(id), 1) FROM promociones));
