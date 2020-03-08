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
        $(`${this.id} #remove`).on("click", (e) => {
            new InterfaceApi({
                instance: this.instance,
                uuids: this.interfaces.store
            }).delete();
        });

        // refresh table and register refresh click.
        $(`${this.id} #refresh`).on("click", (e) => {
            this.table.refresh((e) => {
                this.checkbox.refresh();
            });
        });
        this.table.refresh((e) => {
            this.checkbox.refresh();
        });
    }

    create(data) {
        new InterfaceApi({instance: this.instance}).create(data);
    }

    edit(data) {
        new InterfaceApi({instance: this.instance}).edit(data);
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

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(this.interfaces, e);
            }
        });
        // disabled firstly.
        this.change(this.interfaces, this.interfaces);
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