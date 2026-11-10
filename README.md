# E-commerce Grupo 3

Proyecto con **Arquitectura Monolítica Modular**:
- **Backend:** Go (Golang) — una sola aplicación HTTP organizada por módulos, sirviendo la API REST y opcionalmente los activos web compilados.
- **Frontend:** React + TypeScript con soporte dual **Next.js** y **Vite**.
- **Base de Datos:** PostgreSQL 16 con migraciones SQL y datos semilla.
- **Patrones de Diseño de Software Implementados:**
  - **CRUD 1: Productos — Patrón Composite:** Permite tratar de manera homogénea productos simples (hojas) y combos/packs armados (compuestos) que calculan recursivamente sus precios y evalúan el stock disponible según sus partes.
  - **CRUD 2: Promociones — Patrón Decorator:** Permite envolver dinámicamente un precio base o carrito de compras aplicando capas sucesivas de descuentos porcentuales, montos fijos y cupones con auditoría paso a paso.
- **Despliegue Previsto:** AWS (ECS Fargate, EC2, App Runner o S3 + CloudFront).

---

## Estructura del Proyecto

```text
ecommerce-grupo3-go/
├── backend/
│   ├── cmd/api/             # Punto de entrada HTTP y servidor monolítico
│   ├── internal/
│   │   ├── config/          # Variables de entorno (.env)
│   │   ├── database/        # Conexión PostgreSQL (pgxpool) y auto-migración
│   │   ├── products/        # CRUD 1: Productos (Patrón Composite) + Tests
│   │   └── promotions/      # CRUD 2: Promociones (Patrón Decorator) + Tests
│   ├── go.mod / go.sum
├── database/
│   └── migrations/
│       └── 001_init.sql     # DDL e inserción de datos semilla
├── docker-compose.yml       # Orquestación local inmediata de PostgreSQL
├── frontend/
│   ├── src/
│   │   ├── app/             # Rutas y páginas de Next.js
│   │   ├── services/        # Cliente API REST hacia Go
│   │   ├── types/           # Interfaces TypeScript (Product, Promotion, Steps)
│   │   ├── App.tsx          # Aplicación React interactiva completa
│   │   └── main.tsx         # Punto de entrada Vite
│   ├── next.config.mjs      # Configuración de Next.js
│   ├── vite.config.ts       # Configuración de Vite
│   └── package.json         # Dependencias React 19 + Next + TypeScript
```

---

## 1. Preparar PostgreSQL

### Opción A: Con Docker (Recomendada y más rápida)
Desde la raíz del proyecto ejecuta:
```bash
docker compose up -d
```
Esto levantará PostgreSQL en el puerto `5432` con usuario `postgres`, contraseña `postgrespassword`, base de datos `ecommerce_grupo3` y cargará automáticamente la migración `database/migrations/001_init.sql`.

### Opción B: Con pgAdmin o psql local
1. Crea la base de datos `ecommerce_grupo3`.
2. Ejecuta el script SQL:
```sql
\i database/migrations/001_init.sql
```

---

## 2. Configurar y Ejecutar Backend (Go)

En la carpeta `backend`, asegúrate de que `.env` apunte a tu PostgreSQL:
```env
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgrespassword@localhost:5432/ecommerce_grupo3?sslmode=disable
```

> **Nota de Resiliencia:** Si PostgreSQL aún no está iniciado, el backend Go activa automáticamente un modo en memoria con datos semilla precargados, garantizando que el sistema pueda compilarse y probarse de inmediato sin interrupciones.

Ejecutar el backend:
```bash
cd backend
go run ./cmd/api
```

Ejecutar pruebas unitarias de los patrones y endpoints:
```bash
cd backend
go test -v ./...
```

---

## 3. Ejecutar Frontend (React + Next / Vite)

Abre otra terminal:
```bash
cd frontend
npm install
npm run dev
```

Next.js se iniciará en `http://localhost:3000`.

*(Opcional: Si prefieres usar Vite, ejecuta `npm run vite:dev` en `http://localhost:5173`).*

---

## 4. Detalle de los 2 CRUD y Patrones de Diseño

### CRUD 1: Productos — Patrón Composite
- **Interfaz Component (`Component`):** Implementada por productos individuales y combos.
- **Hoja (`LeafProduct`):** Producto simple con su propio precio y stock.
- **Compuesto (`CompositeProduct`):** Contiene una lista de elementos hijos (`Children []ChildDetail`) y un porcentaje de descuento del combo.
  - `GetPrice()` calcula la suma de todos los componentes hijos y aplica el descuento del paquete.
  - `GetStock()` determina el stock disponible calculando el cuello de botella (mínimo de unidades de componentes disponibles).
- **Rutas API:**
  - `GET /api/productos`: Lista productos simples y combos.
  - `GET /api/productos/{id}`: Consulta un producto individual o combo.
  - `GET /api/productos/{id}/composite`: Desglose detallado del árbol Composite con precios brutos, ahorro y componentes hijos.
  - `POST /api/productos`: Crea un producto simple o un combo indicando sus componentes y descuento.
  - `PUT /api/productos/{id}`: Actualiza un producto o combo.
  - `DELETE /api/productos/{id}`: Elimina un producto o combo.

### CRUD 2: Promociones — Patrón Decorator
- **Interfaz PriceComponent:** Define `GetDescription()`, `GetBasePrice()`, `GetFinalPrice()` y `GetBreakdown()`.
- **Componente Concreto Base (`BasePriceComponent`):** Contiene el precio original de cotización o venta.
- **Decoradores Concretos:**
  - `PercentageDiscountDecorator`: Aplica un descuento porcentual sobre el subtotal acumulado.
  - `FixedAmountDiscountDecorator`: Aplica un bono o rebaja fija en soles (S/).
  - `CouponDiscountDecorator`: Aplica un cupón condicionado a un monto mínimo de compra.
- **Rutas API:**
  - `GET /api/promociones`: Lista todas las promociones.
  - `GET /api/promociones/{id}`: Obtiene el detalle de una promoción.
  - `POST /api/promociones`: Crea una nueva regla promocional.
  - `PUT /api/promociones/{id}`: Actualiza una promoción.
  - `DELETE /api/promociones/{id}`: Elimina una promoción.
  - `POST /api/promociones/aplicar`: Recibe un producto o total de carrito y una lista de decoradores a encadenar, devolviendo la cotización final y la auditoría paso a paso.

---

## 5. Despliegue en AWS

1. **Monolito en Contenedor (AWS ECS Fargate o App Runner):**
   - Se compila el frontend estático (`npm run vite:build` o `next build`).
   - El ejecutable Go sirve tanto la API REST como la aplicación web en el puerto `8080`.
   - Se conecta a una base de datos gestionada **AWS RDS PostgreSQL** mediante la variable de entorno `DATABASE_URL`.
2. **Frontend en S3 + CloudFront y Backend en EC2 / ECS:**
   - La API se expone en un balanceador de carga (ALB).
   - El frontend Next/React se distribuye globalmente con caché CDN de baja latencia.
