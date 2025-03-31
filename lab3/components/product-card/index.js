export class ProductCardComponent {
    constructor(parent) {
        this.parent = parent;
    }

    getHTML(data) {
        return `
            <div class="card" data-id="${data.id}">
                <div class="card-img-container">
                    <img class="card-img-top" src="${data.src}" alt="${data.title}">
                    <div class="card-content-overlay">
                        <h5 class="card-title">${data.title}</h5>
                        <p class="card-text">${data.text}</p>
                    </div>
                </div>
            </div>
        `;
    }
    
    addListeners(listener) {
        const cards = this.parent.querySelectorAll('.card');
        cards.forEach(card => card.addEventListener('click', listener));
    }

    render(data, listener) {
        const html = this.getHTML(data);
        this.parent.insertAdjacentHTML('beforeend', html);
        this.addListeners(listener); 
    }
}