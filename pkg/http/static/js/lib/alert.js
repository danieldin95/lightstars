
export class Alert {

    static div (cls, message, uuid) {
        if (!uuid) {
            uuid = this.genId(32);
        }
        return (`
        <div id="${uuid}" class="alert alert-${cls} alert-dismissible fade show" role="alert">
            ${message}
        </div>`)
    }

    static genId(length) {
        const randomChars = 'abcdefghijklmnopqrstuvwxyz0123456789';
        let result = '';
        for ( let i = 0; i < length; i++ ) {
            result += randomChars.charAt(Math.floor(Math.random() * randomChars.length));
        }
        return result;
    }

    static show(parent, level, message) {
        let uuid = this.genId(32);
        $(parent).append(this.div(level, message, uuid));
        window.setTimeout(function() {
            $(`#${uuid}`).fadeTo('normal', 0).slideUp('normal', function() {
                $(this).remove();
            });
        }, 5000);
    }

    static info (parent, message) {
        this.show(parent, 'info', message)
    }

    static danger (parent, message) {
        this.show(parent, 'danger', message)
    }

    static success (parent, message) {
        this.show(parent, 'success', message);
    }

    static warn (parent, message) {
        this.show(parent, 'warning', message);
    }

    static primary (parent, message) {
        this.show(parent, 'primary', message);
    }
}
