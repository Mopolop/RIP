import { ProductCardComponent } from "../../components/product-card/index.js";
import { ProductPage } from "../product/index.js";
import { productData } from "../../source/data/product-data.js"; // Импортируем данные

export class MainPage {
    constructor(parent) {
        this.parent = parent;
    }
    
    get pageRoot() {
        return document.getElementById('main-page');
    }
        
    getHTML() {
        return `<div id="main-page" class="d-flex flex-wrap"></div>`;
    }
    getData() {
        return [
            { id: 1, src: "https://av.ru/images/hd1/h02/9829392810014.png", title: "Молочные продукты", text: "Молоко, сыр, масло" },
            { id: 2, src: "https://av.ru/product/h46/h03/9745596809246.png", title: "Овощи, фрукты, ягоды", text: "Свежие и замороженные" },
            { id: 3, src: "https://av.ru/images/h7f/h76/9829393399838.png", title: "Мясные продукты", text: "Мясо, колбасы, полуфабрикаты" },
            { id: 4, src: "https://av.ru/images/h9b/h6b/9829393334302.png", title: "Морепродукты", text: "Рыба, креветки, кальмары" },
            { id: 5, src: "https://av.ru/images/hf8/hdc/9829392875550.png", title: "Бакалея", text: "Крупы, макароны, специи" },
            { id: 6, src: "https://av.ru/images/h05/h1c/9829393006622.png", title: "Хлебобулочные изделия", text: "Батоны, булки, пироги" },  
            { id: 7, src: "https://av.ru/product/h2b/h1b/9572362420254.png", title: "Напитки", text: "Соки, вода, газировка" }, 
            { id: 8, src: "source/pictures/Алкоголь.png", title: "Алкоголь", text: "Вина, виски, шампанское" },
        ];
    }

     clickCard(e) {
        const cardElement = e.target.closest('[data-id]');
        if (!cardElement) return;
        
        const categoryId = parseInt(cardElement.dataset.id);
        const productPage = new ProductPage(
            this.parent, 
            categoryId,
            productData // Передаём все данные в ProductPage
        );
        productPage.render();
    }

    render() {
        this.parent.innerHTML = '';
        this.parent.insertAdjacentHTML('beforeend', this.getHTML());

        this.getData().forEach((item) => {
            const productCard = new ProductCardComponent(this.pageRoot);
            productCard.render(item, this.clickCard.bind(this));
        });
    }
}
