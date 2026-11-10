import { Module } from '@nestjs/common';
import { VentasController } from './ventas.controller.js';

@Module({
  controllers: [VentasController]
})
export class VentasModule {}
