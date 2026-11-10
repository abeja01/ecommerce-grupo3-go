import { Injectable } from '@nestjs/common';
import { IVentaService } from './venta.interface.js';

@Injectable()
export class VentaBaseService implements IVentaService {
  procesarVenta(monto: number): string {
    return `Venta procesada exitosamente por un monto de $${monto}`;
  }
}
