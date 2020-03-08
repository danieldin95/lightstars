import {DiskApi} from "./api/disk.js";
import {CheckBoxTop} from "./com/utils.js";
import {DiskTable} from "./widget/disk/table.js";

export class Disk {
    // {
    //   id: '#instance #disk',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        this.id = props.id;
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new Checkbox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new DiskTable({
            id: `${this.id} #display-table`,
            inst: this.inst,
        });

        // register button's click.
        $(`${this.id} #remove`).on("click", (e) => {
            new DiskApi({
                inst: this.inst,
                uuids: this.uuids.store,
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
        new DiskApi({inst: this.inst}).create(data);
    }

    edit(data) {
        new DiskApi({inst: this.inst}).edit(data);
    }
}


export class Checkbox {
    // {
    //   id: '#instane #disk'
    // }
    constructor(props) {
        this.id = props.id;
        this.uuids = {store: [], id: this.id};

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(this.uuids, e);
            },
        });

        // disabled firstly.
        this.change(this.uuids, this.uuids);
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
        } else {
            $(`${record.id} #edit`).removeClass('disabled');
        }
    }
}