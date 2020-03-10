
export class FormWizard {
    constructor(props) {
        this.id = props.id;
        this.default = `${this.id} ${props.default}`;
        this.navtabs = `${this.id} ${props.navigation}`;
        this.form = `${this.id} ${props.form}`;
        this.prev = `${this.id} ${props.buttons.prev}`;
        this.next = `${this.id} ${props.buttons.next}`;
        this.submit = `${this.id} ${props.buttons.submit}`;
        this.cancel = `${this.id} ${props.buttons.cancel}`;
        this.active = this.default;

        this.pages = [];
        $(this.navtabs).each((i, event) => {
            let id = this.id + ' #' + $(event).attr('id');
            this.pages.push(id);
        });

        // register click
        for (let i in this.pages) {
            let page = this.pages[i];
            $(page).on('click', (event) => {
                $(this.active).removeClass('active');
                $($(this.active).attr('data-target')).addClass('d-none');

                this.active = page;
                $(page).addClass('active');
                $($(page).attr('data-target')).removeClass('d-none');
            });
        }
        // reset default page
        this.active = this.default;
        $(this.default).addClass('active');
        $($(this.default).attr('data-target')).removeClass('d-none');

        console.log(this.active, this.pages);
        // register prev and next.
        $(this.prev).on('click', (event) => {
            let pos = this.pages.indexOf(this.active);
            if (pos > 0) {
                let page = this.pages[pos-1];
                $(this.active).removeClass('active');
                $($(this.active).attr('data-target')).addClass('d-none');
                this.active = page;
                $(page).addClass('active');
                $($(page).attr('data-target')).removeClass('d-none');
            }
        });
        $(this.next).on('click', (event) => {
            let pos = this.pages.indexOf(this.active)+1;
            if (pos < this.pages.length) {
                let page = this.pages[pos];
                $(this.active).removeClass('active');
                $($(this.active).attr('data-target')).addClass('d-none');
                this.active = page;
                $(page).addClass('active');
                $($(page).attr('data-target')).removeClass('d-none');
            }
        });
    }

    load (callback) {
        $(this.submit).on('click', (event) => {
            let data = $(this.form).serializeArray();
            if (callback && callback.submit) {
                callback.submit({event, data});
            }
        });
        $(this.cancel).on('click', (event) => {
            if (callback && callback.cancel) {
                callback.cancel({event});
            }
        });
    }
}