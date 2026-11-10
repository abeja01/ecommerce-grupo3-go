import { Controller, Get, Query } from '@nestjs/common';
import { VentaBaseService } from './venta.service.js';
import { VentaLogDecorator } from './venta-log.decorator.js';

@Controller('ventas')
export class VentasController {
  private ventaService = new VentaLogDecorator(new VentaBaseService());

  @Get('procesar')
  procesar(@Query('monto') monto: string) {
    const valorMonto = monto ? parseFloat(monto) : 150.0;
    const resultado = this.ventaService.procesarVenta(valorMonto);

    return {
      mensaje: "Patrón Decorator ejecutado exitosamente en Ventas",
      resultadoServicio: resultado,
      detalle: "Revisa la consola de tu terminal backend para ver los logs de auditoría."
    };
  }
}
