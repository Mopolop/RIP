import { Injectable, Inject, InternalServerErrorException } from '@nestjs/common';
import { CreateStockDto } from './dto/create-stock.dto';
import { UpdateStockDto } from './dto/update-stock.dto';
import { Stock } from './entities/stock.entity';
import { FileService } from '../file/file.service';

@Injectable()
export class StocksService {
  constructor(
    @Inject('STOCKS_FILE_SERVICE') private fileService: FileService<any[]>
  ) { }

  create(createStockDto: CreateStockDto): Stock {
    try {
      const stocks = this.fileService.read();
      const stock = { ...createStockDto, id: stocks.length + 1 };
      this.fileService.add(stock);
      return stock;
    } catch (error) {
      throw new InternalServerErrorException('Failed to create stock');
    }
  }

  findAll(title?: string, text?: string): Stock[] {
    try {
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
    } catch (error) {
      throw new InternalServerErrorException('Failed to fetch stocks');
    }
  }

  findOne(id: number): Stock | null {
    try {
      const stocks = this.fileService.read();
      return stocks.find((stock) => stock.id === id) ?? null;
    } catch (error) {
      throw new InternalServerErrorException('Failed to fetch stock');
    }
  }

  update(id: number, updateStockDto: UpdateStockDto): Stock {
    try {
      const stocks = this.fileService.read();
      const stockIndex = stocks.findIndex((stock) => stock.id === id);

      const updatedStock = { ...stocks[stockIndex], ...updateStockDto };
      stocks[stockIndex] = updatedStock;
      this.fileService.write(stocks);
      return updatedStock;
    } catch (error) {
      throw new InternalServerErrorException('Failed to update stock');
    }
  }

  remove(id: number): boolean {
    try {
      const stocks = this.fileService.read();
      const initialLength = stocks.length;
      const filteredStocks = stocks.filter((stock) => stock.id !== id);

      if (filteredStocks.length === initialLength) {
        return false;
      }

      this.fileService.write(filteredStocks);
      return true;
    } catch (error) {
      throw new InternalServerErrorException('Failed to delete stock');
    }
  }
}