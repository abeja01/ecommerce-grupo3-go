export interface Promotion {
  id: number;
  codigo: string;
  nombre: string;
  descripcion: string;
  tipo: 'porcentaje' | 'monto_fijo' | 'cupon';
  valor: number;
  monto_minimo: number;
  activa: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface PromotionInput {
  codigo: string;
  nombre: string;
  descripcion: string;
  tipo: 'porcentaje' | 'monto_fijo' | 'cupon';
  valor: number;
  monto_minimo: number;
  activa: boolean;
}

export interface PriceStep {
  capa_decorador: string;
  descripcion: string;
  precio_anterior: number;
  descuento_aplicado: number;
  precio_resultante: number;
}

export interface ApplyDecoratorResponse {
  producto_id: number;
  nombre_producto: string;
  precio_base: number;
  descuento_total: number;
  precio_final: number;
  total_decoradores_aplicados: number;
  pasos_decorador: PriceStep[];
}
