
export class Alert {
    static info (message) {
        return this.div('info', message)
    }

    static danger (message) {
        return this.div('danger', message)
    }

    static success (message) {
        return this.div('success', message)
    }

    static warn (message) {
        return this.div('warning', message)
    }

    static primary (message) {
        return this.div('primary', message)
    }

    static div (cls, message) {
        return (`
        <div class="alert alert-${cls} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>`)
    }
}