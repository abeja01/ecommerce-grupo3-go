import { IInventarioComponent } from './item-menu.interface.js';

export class ProductoIndividual implements IInventarioComponent {
  constructor(private nombre: string, private precio: number) {}

  calcularPrecio(): number {
    return this.precio;
  }

  obtenerNombre(): string {
    return this.nombre;
  }
}
