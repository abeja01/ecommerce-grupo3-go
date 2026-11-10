import type { Product, ProductInput, CompositeDetail } from '../types/product';
import type { Promotion, PromotionInput, ApplyDecoratorResponse } from '../types/promotion';

const API_BASE = (typeof window !== 'undefined' && (window as any).__API_URL__) ||
  (typeof process !== 'undefined' && process.env.NEXT_PUBLIC_API_URL) ||
  'http://localhost:8080/api';

// ==========================================
// CRUD 1: PRODUCTOS & COMPOSITE PATTERN
// ==========================================

export async function getProducts(): Promise<Product[]> {
  const res = await fetch(`${API_BASE}/productos`);
  if (!res.ok) throw new Error('Error al listar productos');
  return res.json();
}

export async function getProduct(id: number): Promise<Product> {
  const res = await fetch(`${API_BASE}/productos/${id}`);
  if (!res.ok) throw new Error(`Error al consultar producto #${id}`);
  return res.json();
}

export async function getCompositeDetail(id: number): Promise<CompositeDetail> {
  const res = await fetch(`${API_BASE}/productos/${id}/composite`);
  if (!res.ok) throw new Error('Error al obtener estructura Composite');
  return res.json();
}

export async function createProduct(input: ProductInput): Promise<Product> {
  const res = await fetch(`${API_BASE}/productos`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Error al crear producto');
  }
  return res.json();
}

export async function updateProduct(id: number, input: ProductInput): Promise<Product> {
  const res = await fetch(`${API_BASE}/productos/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Error al actualizar producto');
  }
  return res.json();
}

export async function deleteProduct(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/productos/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Error al eliminar producto');
}

// ==========================================
// CRUD 2: PROMOCIONES & DECORATOR PATTERN
// ==========================================

export async function getPromotions(): Promise<Promotion[]> {
  const res = await fetch(`${API_BASE}/promociones`);
  if (!res.ok) throw new Error('Error al listar promociones');
  return res.json();
}

export async function getPromotion(id: number): Promise<Promotion> {
  const res = await fetch(`${API_BASE}/promociones/${id}`);
  if (!res.ok) throw new Error(`Error al consultar promoción #${id}`);
  return res.json();
}

export async function createPromotion(input: PromotionInput): Promise<Promotion> {
  const res = await fetch(`${API_BASE}/promociones`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Error al crear promoción');
  }
  return res.json();
}

export async function updatePromotion(id: number, input: PromotionInput): Promise<Promotion> {
  const res = await fetch(`${API_BASE}/promociones/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Error al actualizar promoción');
  }
  return res.json();
}

export async function deletePromotion(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/promociones/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Error al eliminar promoción');
}

export async function applyPromotionDecorators(req: {
  producto_id?: number;
  monto_personalizado?: number;
  promocion_ids?: number[];
  codigo_cupon?: string;
}): Promise<ApplyDecoratorResponse> {
  const res = await fetch(`${API_BASE}/promociones/aplicar`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || 'Error al aplicar decoradores');
  }
  return res.json();
}

export async function getHealth(): Promise<any> {
  const res = await fetch(`${API_BASE}/health`);
  return res.json();
}
