export class ProductCardComponent {
    constructor(parent) {
        this.parent = parent;
    }

    getHTML(data) {
        return `
            <div class="card" style="width: 300px; cursor: pointer;" data-id="${data.id}">
                <img class="card-img-top" src="${data.src}" alt="картинка">
                <div class="card-body">
                    <h5 class="card-title">${data.title}</h5>
                    <p class="card-text">${data.text}</p>
                </div>  
            </div>
        `;
    }
    
    addListeners(data, listener) {
        // Добавляем обработчик на всю карточку
        this.parent.lastElementChild.addEventListener("click", listener);
    }
    
    render(data, listener) {
        const html = this.getHTML(data);
        this.parent.insertAdjacentHTML('beforeend', html);
        this.addListeners(data, listener);
    }
}