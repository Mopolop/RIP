import { ajax } from "../../modules/ajax.js";
import { stockUrls } from "../../modules/stockUrls.js";
import { ProductPage } from "../product/index.js";

export class AddProductPage {
    constructor(parent, categoryId) {
        this.parent = parent;
        this.categoryId = categoryId;
    }

    getHTML() {
        return `
            <div class="add-product-form">
                <h2>Добавить товар</h2>
                <form id="productForm">
                    <input type="text" name="title" placeholder="Название" required><br>
                    <input type="text" name="src" placeholder="Ссылка на картинку" required><br>
                    <input type="text" name="text" placeholder="Описание" required><br>
                    <input type="text" name="badgeText" placeholder="Текст бейджа (необязательно)"><br>
                    <select name="badgeClass" class="form-select" style="margin-bottom: 18px;">
                        <option value="">Выберите стиль бейджа (необязательно)</option>
                        <option value="badge text-bg-primary">Синий (primary)</option>
                        <option value="badge text-bg-secondary">Серый (secondary)</option>
                        <option value="badge text-bg-success">Зелёный (success)</option>
                        <option value="badge text-bg-danger">Красный (danger)</option>
                        <option value="badge text-bg-warning">Жёлтый (warning)</option>
                        <option value="badge text-bg-info">Голубой (info)</option>
                        <option value="badge text-bg-light">Светлый (light)</option>
                        <option value="badge text-bg-dark">Тёмный (dark)</option>
                    </select><br>
                    <button type="submit" class="btn btn-success">Создать</button>
                    <button type="button" class="btn btn-secondary back-btn">Назад</button>
                </form>
            </div>
        `;
    }

    render() {
        // Проверяем, не подключён ли уже стиль
        if (!document.getElementById('add-product-css')) {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = 'source/styles/add-product.css';
            link.id = 'add-product-css';
            document.head.appendChild(link);
        }

        this.parent.innerHTML = this.getHTML();
        document.getElementById('productForm').onsubmit = this.submitForm.bind(this);
        this.parent.querySelector('.back-btn').onclick = () => {
            new ProductPage(this.parent, this.categoryId).render();
        };
    }

    submitForm(e) {
        e.preventDefault();
        const form = e.target;
        const data = {
            title: form.title.value,
            src: form.src.value,
            text: form.text.value
        };
        // Добавляем бейдж, если заполнено
        if (form.badgeText.value && form.badgeClass.value) {
            data.badge = {
                text: form.badgeText.value,
                class: form.badgeClass.value
            };
        }
        ajax.post(`${stockUrls.getStocks()}/${this.categoryId}/products`, data, () => {
            new ProductPage(this.parent, this.categoryId).render();
        });
    }
}