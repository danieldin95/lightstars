export var Alert = {
    info: function (message) {
        return this.div('info', message)
    },
    danger: function (message) {
        return this.div('danger', message)
    },
    success: function (message) {
        return this.div('success', message)
    },
    warn: function (message) {
        return this.div('warning', message)
    },
    primary: function (message) {
        return this.div('primary', message)
    },
    div: function (cls, message) {
        return (`
        <div class="alert alert-${cls} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>`)
    },
};