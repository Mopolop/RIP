export class ProductComponent {
    constructor(parent) {
        this.parent = parent;
    }

    getHTML(data) {
        // Генерируем бейдж, если есть данные
        const badge = data.badge 
            ? `<span class="badge ${data.badge.class}">${data.badge.text}</span>`
            : '<br>';

        return `
            <div class="product-card">
                <img src="${data.src}" 
                     alt="${data.title}" 
                     class="product-card-image">
                <div class="product-card-body">
                    <div class="card-header">
                        ${badge}
                        <h3 class="product-card-title">${data.title}</h3>    
                    </div>
                    <p class="product-card-text">${data.text}</p>
                </div>
            </div>
        `;
    }

    render(data) {
        const html = this.getHTML(data);
        this.parent.insertAdjacentHTML('beforeend', html);
    }
}