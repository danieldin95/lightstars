
export class FormWizard {
    //
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
                this.move(i);
            });
        }
        // reset default page
        this.active = this.default;
        $(this.default).addClass('active');
        $(this.target(this.default)).removeClass('d-none');

        console.log(this.active, this.pages);
        // register prev and next.
        $(this.prev).on('click', (event) => {
            let pos = this.pages.indexOf(this.active);
            if (pos > 0) {
                this.move(pos-1);
            }
        });
        $(this.next).on('click', (event) => {
            let pos = this.pages.indexOf(this.active)+1;
            if (pos < this.pages.length) {
                this.move(pos);
            }
        });
    }

    target(page) {
        return this.id + " " + $(page).attr('data-target');
    }

    move (pos) {
        let page = this.pages[pos];
        $(this.active).removeClass('active');
        $(this.target(this.active)).addClass('d-none');
        this.active = page;
        $(page).addClass('active');
        $(this.target(page)).removeClass('d-none');
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