export interface ComboItem {
  product_id: number;
  nombre?: string;
  precio_unitario?: number;
  cantidad: number;
  subtotal?: number;
}

export interface Product {
  id: number;
  category_id?: number;
  nombre: string;
  descripcion: string;
  precio: number;
  stock: number;
  imagen_url: string;
  tipo: 'simple' | 'combo';
  descuento_porcentaje: number;
  activo: boolean;
  created_at?: string;
  items?: ComboItem[];
}

export interface ProductInput {
  category_id?: number;
  nombre: string;
  descripcion: string;
  precio: number;
  stock: number;
  imagen_url: string;
  tipo: 'simple' | 'combo';
  descuento_porcentaje: number;
  activo: boolean;
  items: ComboItem[];
}

export interface CompositeDetail {
  id: number;
  nombre: string;
  tipo: string;
  precio_final: number;
  precio_bruto: number;
  ahorro: number;
  descuento_porcentaje: number;
  stock_disponible: number;
  elementos_hijos: ComboItem[];
}
