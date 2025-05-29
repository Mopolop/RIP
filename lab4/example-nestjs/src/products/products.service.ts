import { Injectable, Inject } from '@nestjs/common';
import { FileService } from '../file/file.service';

@Injectable()
export class ProductsService {
  constructor(
    @Inject('PRODUCTS_FILE_SERVICE') private fileService: FileService<Record<string, any[]>>
  ) {}

  getProductsByCategory(categoryId: string): any[] {
    const productsData = this.fileService.read();
    return productsData[categoryId] || [];
  }

 addProductToCategory(categoryId: string, product: any): any {
    const productsData = this.fileService.read();
    const categoryProducts = productsData[categoryId] || [];
    
    // Генерация ID
    const nextNumber = categoryProducts.length + 1;
    product.id = parseInt(`${categoryId}${nextNumber}`);

    // Обработка badge
    if (product.badgeText && product.badgeClass) {
      product.badge = {
        text: product.badgeText,
        class: product.badgeClass
      };
      delete product.badgeText;
      delete product.badgeClass;
    }

    categoryProducts.push(product);
    productsData[categoryId] = categoryProducts;
    this.fileService.write(productsData);

    return product;
  }
  
  deleteProduct(categoryId: string, productId: string): boolean {
    const productsData = this.fileService.read();
    if (!productsData[categoryId]) {
      return false;
    }
    const initialLength = productsData[categoryId].length;
    productsData[categoryId] = productsData[categoryId].filter(
      p => p.id != productId
    );
    if (productsData[categoryId].length === initialLength) {
      return false;
    }
    this.fileService.write(productsData);
    return true;
  }
}