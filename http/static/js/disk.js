import {DiskApi} from "./api/disk.js";
import {DiskTable} from "./widget/disk/table.js";
import {CheckBoxTab} from "./widget/checkbox/checkbox.js";


class CheckBox extends CheckBoxTab {
}


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

        this.checkbox = new CheckBox(props);
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