// stocks.service.ts
import { Injectable } from '@nestjs/common';
import { CreateStockDto } from './dto/create-stock.dto';
import { UpdateStockDto } from './dto/update-stock.dto';
import { Stock } from './entities/stock.entity';
import { FileService } from 'src/file/file.service';

@Injectable()
export class StocksService {
  constructor(private fileService: FileService<Stock[]>) {}

  create(createStockDto: CreateStockDto) {
    const stocks = this.fileService.read();
    const stock = { ...createStockDto, id: stocks.length + 1 };
    this.fileService.add(stock);
  }

  findAll(title?: string, text?: string): Stock[] {
    const stocks = this.fileService.read();

    let filteredStocks = stocks;

    if (title) {
      filteredStocks = filteredStocks.filter((stock) =>
        stock.title.toLowerCase().includes(title.toLowerCase())
      );
    }

    if (text) {
      filteredStocks = filteredStocks.filter((stock) =>
        stock.text.toLowerCase().includes(text.toLowerCase())
      );
    }

    return filteredStocks;
  }

  findOne(id: number): Stock | null {
    const stocks = this.fileService.read();
    return stocks.find((stock) => stock.id === id) ?? null;
  }

  update(id: number, updateStockDto: UpdateStockDto): void {
    const stocks = this.fileService.read();
    const updatedStocks = stocks.map((stock) =>
      stock.id === id ? { ...stock, ...updateStockDto } : stock,
    );
    this.fileService.write(updatedStocks);
  }

  remove(id: number): void {
    const filteredStocks = this.fileService
      .read()
      .filter((stock) => stock.id !== id);
    this.fileService.write(filteredStocks);
  }
}