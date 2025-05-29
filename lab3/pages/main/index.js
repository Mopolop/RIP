import { ProductCardComponent } from "../../components/product-card/index.js";
import { ProductPage } from "../product/index.js";
import {ajax} from "../../modules/ajax.js";
import {stockUrls} from "../../modules/stockUrls.js";

export class MainPage {
    constructor(parent) {
        this.parent = parent;
    }
    
    get pageRoot() {
        return document.getElementById('main-page');
    }
        
    getHTML() {
        return `
            <div class="header-container">
                <div class="header-rectangle">
                    <h1>Витязь-Аэро</h1>
                    <p>Широкий ассортимент продуктов на любой вкус</p>
                </div>
            </div>
            <div id="main-page" class="d-flex flex-wrap"></div>
        `;
    }

getData() {
    ajax.get(stockUrls.getStocks(), (data) => {
        this.renderData(data);
    })
}
 renderData(items) {
    items.forEach((item) => {
        const productCard = new ProductCardComponent(this.pageRoot)
        productCard.render(item, this.clickCard.bind(this))
    })
}

    clickCard(e) {
        const cardElement = e.target.closest('[data-id]');
        if (!cardElement) return;
        
        const categoryId = parseInt(cardElement.dataset.id);
        const productPage = new ProductPage(
            this.parent, 
            categoryId,
            null // Данные будут загружаться через API
        );
        productPage.render();
    }

    render() {
        this.parent.innerHTML = '';
        this.parent.insertAdjacentHTML('beforeend', this.getHTML());
        this.getData();
    }
}
