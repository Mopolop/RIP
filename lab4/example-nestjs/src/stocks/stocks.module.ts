import { Module } from '@nestjs/common';
import { StocksController } from './stocks.controller';
import { StocksService } from './stocks.service';
import { FileService } from '../file/file.service';
import { Stock } from './entities/stock.entity';

@Module({
  controllers: [StocksController],
  providers: [
    StocksService,
    {
      provide: FileService,
      useFactory: () => new FileService<Stock[]>('stocks.json'),
    },
  ],
})
export class StocksModule {}