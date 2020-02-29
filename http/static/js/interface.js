import {InterfaceApi} from "./api/interface.js";
import {CheckBoxTop} from "./com/utils.js";
import {InterfaceTable} from "./widget/interface/table.js";

export class Interface {
    // {
    //   id: '#instance #interface',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.name = props.name;
        this.instance = props.uuid;

        this.checkbox = new Checkbox(props);
        this.interfaces = this.checkbox.interfaces;
        this.table = new InterfaceTable({
            id: `${this.id} #display-table`,
            instance: this.instance,
        });

        // register buttons's click
        $(`${this.id} #remove`).on("click", this, function (e) {
            new InterfaceApi({
                instance: e.data.instance,
                uuids: e.data.interfaces.store,
                name: this.name}).delete();
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
        new InterfaceApi({instance: this.instance, name: this.name}).create(data);
    }

    edit(data) {
        new InterfaceApi({instance: this.instance, name: this.name}).edit(data);
    }
}


export class Checkbox {
    // {
    //   id: '#instane #disk'
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.interfaces = {store: [], id: this.id};

        let record = this.interfaces;
        let change = this.change;

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: function (e) {
                change(record, e);
            }
        });
        // disabled firstly.
        change(record, this.interfaces);
    }

    refresh() {
        this.top.refresh();
    }

    change(record, from) {
        record.store = from.store;

        if (from.store.length == 0) {
            $(`${record.id} #remove`).addClass('disabled');
        } else {
            $(`${record.id} #remove`).removeClass('disabled');
        }
        if (from.store.length != 1) {
            $(`${record.id} #edit`).addClass('disabled');
        }  else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}