import {Ctl} from "./ctl.js";
import {DiskApi} from "../api/disk.js";
import {DiskTable} from "../widget/disk/table.js";
import {CheckBox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckBox {
}


export class DiskCtl extends Ctl {
    // {
    //   id: '#instance #disk',
    //   uuid: uuid of instance,
    //   name: name of instance,
    //   onRemove: function {},
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new DiskTable({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", (e) => {
            let data = {inst: this.inst, uuids: this.uuids.store};
            if (this.props.onRemove) {
                this.props.onRemove(data);
            } else {
                new DiskApi(data).delete();
            }
        });

        // refresh table and register refresh click.
        $(this.child('#refresh')).on("click", (e) => {
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
