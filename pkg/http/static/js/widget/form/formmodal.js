//
export class FormModalWid {
    // {
    //   id: '#disk'
    // }
    constructor(props) {
        this.props = props;
        this.id = props.id;
        this.forms = `${this.id} form`;

        this.events = {
            submit: {
                func: function (e) {
                },
                data: undefined,
            }
        };
    }

    template() {
        return (`<not-implement/>`);
    }

    render() {
        this.view = $(this.template());
        this.view.find('label[for]').each(function () {
            let label = $(this);
            let targetId = (label.attr('for') || '').trim();
            if (!targetId) {
                return;
            }
            if (label.closest('.modal-content').find(`#${targetId}`).length > 0) {
                return;
            }
            let ctl = label.closest('.form-group').find('input,select,textarea').first();
            if (ctl.length > 0) {
                ctl.attr('id', targetId);
            }
        });
        this.view.find('input:not([autocomplete])').each(function () {
            let type = (this.type || "").toLowerCase();
            if (type === "password") {
                $(this).attr('autocomplete', 'current-password');
            } else {
                $(this).attr('autocomplete', 'off');
            }
        });
        this.view.find('textarea:not([autocomplete])').attr('autocomplete', 'off');
        this.container().html(this.view);
    }

    fetch() {
        console.log('not-implement')
    }

    submit() {
        if (this.events.submit.func) {
            this.events.submit.func({
                data: this.events.submit.data,
                form: $(this.forms).serializeArray(),
            });
        }
    }

    container() {
        return $(this.id);
    }

    onsubmit(data, func) {
        if (typeof data == "function") {
            this.events.submit.data = {};
            this.events.submit.func = data;
        } else {
            this.events.submit.data = data;
            this.events.submit.func = func;
        }
    }

    loading() {
        let releaseFocus = (container) => {
            let root = container && container.length ? container[0] : null;
            let active = document.activeElement;
            if (active && root && root.contains(active) && typeof active.blur === "function") {
                active.blur();
            }
            if (root && typeof root.blur === "function") {
                root.blur();
            }
        };

        this.container().off('hidden.bs.modal.formmodal');
        this.container().on('hidden.bs.modal.formmodal', function() {
            // Move focus out of hidden modal for a11y.
            if (document.activeElement && typeof document.activeElement.blur === "function") {
                document.activeElement.blur();
            }
            if (document.body && typeof document.body.focus === "function") {
                document.body.focus();
            }
        });

        this.container().off('click.formmodal', '[name=finish-btn]');
        this.container().on('click.formmodal', '[name=finish-btn]', this, function(e) {
            e.preventDefault();
            e.stopPropagation();
            releaseFocus(e.data.container());
            e.data.submit();
            e.data.container().modal("hide");
        });
        this.container().off('click.formmodal', '[name=cancel-btn]');
        this.container().on('click.formmodal', '[name=cancel-btn]', this, function(e) {
            e.preventDefault();
            e.stopPropagation();
            releaseFocus(e.data.container());
            e.data.container().modal("hide");
        });
        $(this.forms).each(function (i, e) {
            $(e).on('submit', function (e) {
                return false;
            });
        });
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }
}
