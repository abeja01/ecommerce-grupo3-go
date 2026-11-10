import { Controller, Get } from '@nestjs/common';
import { ProductoIndividual } from './producto-individual.js';
import { CategoriaCombo } from './categoria-combo.js';

@Controller('inventario')
export class InventarioController {
  @Get('prueba-composite')
  probarComposite() {
    const laptop = new ProductoIndividual("Laptop Gamer", 1200);
    const mouse = new ProductoIndividual("Mouse Inalámbrico", 30);
    const teclado = new ProductoIndividual("Teclado Mecánico", 80);

    const comboGamer = new CategoriaCombo("Pack Gamer Pro");
    comboGamer.agregarItem(laptop);
    comboGamer.agregarItem(mouse);
    comboGamer.agregarItem(teclado);

    return {
      mensaje: "Patrón Composite ejecutado exitosamente en el Inventario",
      nombreCombo: comboGamer.obtenerNombre(),
      precioTotalCalculado: comboGamer.calcularPrecio()
    };
  }
}
