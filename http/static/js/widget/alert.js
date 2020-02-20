export default class Alert {
    constructor(message) {
        this.message = message;
    }
    danger() {
        return `
            <div class="alert alert-danger alert-dismissible fade show" role="alert">
                ${this.message}
                <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>`
    }
    warning() {
        return `
            <div class="alert alert-warning alert-dismissible fade show" role="alert">
                ${this.message}
                <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>`
    }
    info() {
        return `
            <div class="alert alert-info alert-dismissible fade show" role="alert">
                ${this.message}
                <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>`
    }
}