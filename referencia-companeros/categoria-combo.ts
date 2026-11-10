import { IInventarioComponent } from './item-menu.interface.js';

export class CategoriaCombo implements IInventarioComponent {
  private items: IInventarioComponent[] = [];

  constructor(private nombre: string) {}

  agregarItem(item: IInventarioComponent): void {
    this.items.push(item);
  }

  calcularPrecio(): number {
    let total = 0;
    for (const item of this.items) {
      total += item.calcularPrecio();
    }
    return total;
  }

  obtenerNombre(): string {
    return this.nombre;
  }
}
