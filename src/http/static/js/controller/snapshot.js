import {Ctl} from "./ctl.js";
import {SnapshotApi} from "../api/snapshot.js";
import {SnapshotTable} from "../widget/snapshot/table.js";
import {CheckBox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckBox {
    change(from) {
        super.change(from);
        if (from.store.length === 1) {
            $(this.child('#revert')).removeAttr('disabled');
        } else {
            $(this.child('#revert')).attr("disabled","disabled");
        }
    }
}


export class SnapshotCtl extends Ctl {
    // {
    //   id: '#instance #snapshot',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new SnapshotTable({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", this, function (e) {
            new SnapshotApi({
                inst: e.data.inst,
                uuids: e.data.uuids.store,
                name: e.data.name}).delete();
        });

        $(this.child('#revert')).on("click", this, function (e) {
            new SnapshotApi({
                inst: e.data.inst,
                uuids: e.data.uuids.store,
                name: e.data.name}).revert();
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
        new SnapshotApi({inst: this.inst, name: this.name}).create(data);
    }
}
