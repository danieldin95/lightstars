import {CheckBoxAll} from "../../com/checkbox.js";


export class CheckBoxTab {
    // {
    //   id: "#networks"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.uuids = {store: [], id: this.id};

        this.top = new CheckBoxAll({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(e);
            },
        });

        // disable firstly.
        this.change(this.uuids);
    }

    refresh() {
        this.top.refresh();
    }

    change(from) {
        this.uuids.store = from.store;

        if (from.store.length === 0) {
            $(`${this.uuids.id} #delete`).addClass('disabled');
        } else {
            $(`${this.uuids.id} #delete`).removeClass('disabled');
        }
        if (from.store.length !== 1) {
            $(`${this.uuids.id} #edit`).addClass('disabled');
        } else {
            $(`${this.uuids.id} #edit`).removeClass('disabled');
        }
    }
}