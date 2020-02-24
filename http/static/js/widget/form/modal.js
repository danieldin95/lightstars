
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
                array: $(this.forms).serializeArray()
            });
        }
    }

    container() {
        return $(`#${this.containerId}`);
    }

    onSubmit(data, fn) {
        this.events.submit.data = data;
        this.events.submit.func = fn;
    }
}