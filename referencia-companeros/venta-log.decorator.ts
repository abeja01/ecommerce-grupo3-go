import { IVentaService } from './venta.interface.js';

export class VentaLogDecorator implements IVentaService {
  constructor(private ventaService: IVentaService) {}

  procesarVenta(monto: number): string {
    console.log(`[LOG DE AUDITORÍA]: Iniciando venta por $${monto} a las ${new Date().toISOString()}`);
    const resultado = this.ventaService.procesarVenta(monto);
    console.log(`[LOG DE AUDITORÍA]: Venta finalizada con éxito.`);
    return resultado;
  }
}
