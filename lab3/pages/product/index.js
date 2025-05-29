import { BackButtonComponent } from "../../components/back-button/index.js";
import { MainPage } from "../main/index.js";
import { ProductComponent } from "../../components/product/index.js";
import { ajax } from "../../modules/ajax.js";
import { stockUrls } from "../../modules/stockUrls.js";
import { AddProductPage } from "../add-product/index.js";

export class ProductPage {
    constructor(parent, categoryId) {
        this.parent = parent;
        this.categoryId = categoryId;
    }

    get pageRoot() {
        return document.getElementById("product-page");
    }

    getHTML() {
        return `
            <div id="product-page">
                <div class="button-container">
                    <button class="action-button" id="back-btn">Назад</button>
                    <button class="action-button" id="add-product-btn">Добавить товар</button>
                </div>
                <div class="products-container"></div>
            </div>
        `;
    }

    getData() {
        ajax.get(`${stockUrls.getStocks()}/${this.categoryId}/products`, (data) => {
            this.renderData(data);
        });
    }

    renderData(items) {
        const container = this.pageRoot.querySelector('.products-container');
        container.innerHTML = '';
        items.forEach(item => {
            const product = new ProductComponent(container);
            product.render(item);
        });
    }

    clickBack() {
        const mainPage = new MainPage(this.parent);
        mainPage.render();
    }

    render() {
        this.parent.innerHTML = '';
        const html = this.getHTML();
        this.parent.insertAdjacentHTML('beforeend', html);

        const backButton = this.pageRoot.querySelector('#back-btn');
        backButton.addEventListener('click', this.clickBack.bind(this));

        const addButton = this.pageRoot.querySelector('#add-product-btn');
        addButton.addEventListener('click', () => {
            const addProductPage = new AddProductPage(this.parent, this.categoryId);
            addProductPage.render();
        });

        this.getData();
    }
}
