import { BackButtonComponent } from "../../components/back-button/index.js";
import { MainPage } from "../main/index.js";
import { ProductComponent } from "../../components/product/index.js";

export class ProductPage {
    constructor(parent, categoryId, productData) {
        this.parent = parent;
        this.categoryId = categoryId; // ID категории (1-8)
        this.productData = productData; // Все данные о товарах
    }

    getProductsForCategory() {
        // Получаем товары для выбранной категории
        return this.productData[this.categoryId] || [];
    }
    
    get pageRoot() {
        return document.getElementById("product-page");
    }

    getHTML() {
        return `
            <div id="product-page">
                <div class="products-container"></div>
            </div>
        `;
    }

    
    clickBack() {
        const mainPage = new MainPage(this.parent);
        mainPage.render();
    }

    render() {
        this.parent.innerHTML = "";
        this.parent.insertAdjacentHTML("beforeend", this.getHTML());
        
        // Кнопка "Назад"
        const backButton = new BackButtonComponent(this.pageRoot);
        backButton.render(this.clickBack.bind(this));

        // Получаем товары для категории
        const products = this.getProductsForCategory();
        
        if (products.length === 0) {
            this.pageRoot.innerHTML += "<p>Товары не найдены</p>";
            return;
        }

        // Рендерим все товары категории
        const container = this.pageRoot.querySelector('.products-container');
        products.forEach(product => {
            const productComponent = new ProductComponent(container);
            productComponent.render(product);
        });
    }
}
