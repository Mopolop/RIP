import { Module } from '@nestjs/common';
import { StocksController } from './stocks.controller';
import { StocksService } from './stocks.service';
import { ProductsService } from '../products/products.service';
import { FileService } from '../file/file.service';

@Module({
  controllers: [StocksController],
  providers: [
    StocksService,
    ProductsService,
    {
      provide: 'STOCKS_FILE_SERVICE',
      useFactory: () => new FileService<any[]>('stocks.json'),
    },
    {
      provide: 'PRODUCTS_FILE_SERVICE',
      useFactory: () => new FileService<Record<string, any[]>>('product-data.json'),
    },
  ],
})
export class StocksModule {}