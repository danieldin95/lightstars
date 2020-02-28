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
        this.checkbox = new NetworkOn(props);
        this.networks = this.checkbox.uuids;
        this.table = new NetworkTable({id: `${this.id} #display-table`});

        // register buttons's click.
        $(`${this.id} #delete`).on("click", this.networks, function (e) {
            new NetworkApi({uuids: e.data.store}).delete();
        });

        // refresh table and register refresh click.
        let refresh = function (your) {
            your.table.refresh(your.checkbox, function (e) {
                e.data.refresh();
            });
        };
        $(`${this.id} #refresh`).on("click", this, function (e) {
            refresh(e.data);
        });
        refresh(this);
    }

    create(data) {
        new NetworkApi().create(data);
    }
}


export class NetworkOn {
    // {
    //   id: "#networks"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.uuids = {store: [], id: this.id};

        let change = this.change;
        let record = this.uuids;

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: function (e) {
                change(record, e);
            },
        });

        // disabled firstly.
        change(record, this.uuids);
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