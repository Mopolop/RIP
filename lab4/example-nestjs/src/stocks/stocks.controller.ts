import { Controller, Get, Post, Body, Patch, Param, Delete, Query, NotFoundException } from '@nestjs/common';
import { StocksService } from './stocks.service';
import { ProductsService } from '../products/products.service';
import { CreateStockDto } from './dto/create-stock.dto';
import { UpdateStockDto } from './dto/update-stock.dto';

@Controller('stocks')
export class StocksController {
  constructor(
    private readonly stocksService: StocksService,
    private readonly productsService: ProductsService
  ) {}

  @Post()
  create(@Body() createStockDto: CreateStockDto) {
    return this.stocksService.create(createStockDto);
  }

  @Get()
  findAll(
    @Query('title') title?: string,
    @Query('text') text?: string,
    @Query('id') id?: string
  ) {
    if (id) {
      const stock = this.stocksService.findOne(+id);
      return stock ? [stock] : [];
    }
    return this.stocksService.findAll(title, text);
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    const stock = this.stocksService.findOne(+id);
    if (!stock) {
      throw new NotFoundException(`Stock with ID ${id} not found`);
    }
    return stock;
  }
  
  @Get(':id/products')
  getProductsByCategory(@Param('id') id: string) {
    const products = this.productsService.getProductsByCategory(id);
    if (!products) {
      throw new NotFoundException(`Products for category ${id} not found`);
    }
    return products;
  }

  @Post(':id/products')
  addProductToCategory(@Param('id') id: string, @Body() product: any) {
    return this.productsService.addProductToCategory(id, product);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() updateStockDto: UpdateStockDto) {
    const stock = this.stocksService.findOne(+id);
    if (!stock) {
      throw new NotFoundException(`Stock with ID ${id} not found`);
    }
    return this.stocksService.update(+id, updateStockDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    const stock = this.stocksService.findOne(+id);
    if (!stock) {
      throw new NotFoundException(`Stock with ID ${id} not found`);
    }
    return this.stocksService.remove(+id);
  }

  @Delete(':categoryId/products/:productId')
  deleteProduct(
    @Param('categoryId') categoryId: string,
    @Param('productId') productId: string
  ) {
    const success = this.productsService.deleteProduct(categoryId, productId);
    if (!success) {
      throw new NotFoundException('Product or category not found');
    }
    return { success: true };
  }
}