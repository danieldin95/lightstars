import {NetworkApi} from "./api/network.js";
import {CheckBoxTop} from "./com/utils.js";
import {NetworkTable} from "./widget/network/table.js";


export class Network {
    // {
    //   id: "#networks"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.checkbox = new Checkbox(props);
        this.networks = this.checkbox.uuids;
        this.table = new NetworkTable({id: `${this.id} #display-table`});

        // register buttons's click.
        $(`${this.id} #delete`).on("click", (e) => {
            new NetworkApi({uuids: this.networks.store}).delete();
        });

        // refresh table and register refresh click.
        $(`${this.id} #refresh`).on("click", (e) => {
            this.table.refresh((e) => {
                this.checkbox.refresh();
            })
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });
    }

    create(data) {
        new NetworkApi().create(data);
    }
}


export class Checkbox {
    // {
    //   id: "#networks"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.uuids = {store: [], id: this.id};

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(this.uuids, e);
            },
        });

        // disable firstly.
        this.change(this.uuids, this.uuids);
    }

    refresh() {
        this.top.refresh();
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length == 0) {
            $(`${record.id} #edit`).addClass('disabled');
            $(`${record.id} #delete`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
            $(`${record.id} #delete`).removeClass('disabled');
        }
        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}