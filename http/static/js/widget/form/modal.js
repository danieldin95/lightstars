

export class ModalFormBase {
    // {containerId: ''}
    constructor(props) {
        this.props = props;
        this.containerId = props.containerId;
        this.forms = `#${this.containerId} form`;
        this.events = {
            submit: {
                func: function (e) {
                },
                data: undefined,
            }
        };
    }

    template() {
        return `<NotImplate></NotImplate>`
    }

    render() {
        this.view = $(this.template());
        this.container().html(this.view);
    }

    fetch() {
        console.log('NotImplement')
    }

    submit() {
        if (this.events.submit.func) {
            this.events.submit.func({
                data: this.events.submit.data,
                array: $(this.forms).serializeArray(),
            });
        }
        return false;
    }

    container() {
        return $(`#${this.containerId}`);
    }

    onsubmit(data, fn) {
        this.events.submit.data = data;
        this.events.submit.func = fn;
    }

    loading() {
        this.container().find('[name=finish-btn]').on('click', this, function(e) {
            e.data.submit();
            e.data.container().modal("hide");
        });
        this.container().find('[name=cancel-btn]').on('click', this, function(e) {
            e.data.container().modal("hide");
        });

        $(this.forms).each(function (i, e) {
            console.log("disable form's submit", e);
            $(e).on('submit', function (e) {
                return false;
            });
        });
    }
}