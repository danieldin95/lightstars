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
        this.instance = props.uuid;

        this.checkbox = new Checkbox(props);
        this.disks = this.checkbox.disks;
        this.table = new DiskTable({
            id: `${this.id} #display-table`,
            instance: this.instance,
        });

        // register button's click.
        $(`${this.id} #remove`).on("click", (e) => {
            new DiskApi({
                instance: this.instance,
                uuids: this.disks.store,
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
        new DiskApi({instance: this.instance}).create(data);
    }

    edit(data) {
        new DiskApi({instance: this.instance}).edit(data);
    }
}


export class Checkbox {
    // {
    //   id: '#instane #disk'
    // }
    constructor(props) {
        this.id = props.id;
        this.disks = {store: [], id: this.id};

        this.top = new CheckBoxTop({
            one: `${this.id} #on-one`,
            all: `${this.id} #on-all`,
            change: (e) => {
                this.change(this.disks, e);
            },
        });

        // disabled firstly.
        this.change(this.disks, this.disks);
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