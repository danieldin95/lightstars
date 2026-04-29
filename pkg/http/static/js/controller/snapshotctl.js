import {Controller} from "./controller.js";
import {SnapshotApi} from "../api/snapshotapi.js";
import {SnapshotTableWid} from "../widget/snapshot/snapshottable.js";
import {CheckboxWid} from "../widget/common/checkbox.js";
import {ConfirmActionWid} from "../widget/common/confirmaction.js";


class CheckBoxCtl extends CheckboxWid {
    change(from) {
        super.change(from);
        if (from.store.length === 1) {
            $(this.child('#revert')).removeAttr('disabled');
        } else {
            $(this.child('#revert')).attr("disabled","disabled");
        }
    }
}


export class SnapshotCtl extends Controller {
    // {
    //   id: '#instance #snapshot',
    //   uuid: uuid of instance,
    //   name: name of instance,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.inst = props.uuid;
        this.confirm = props.confirm;

        this.CheckboxWid = new CheckBoxCtl(props);
        this.uuids = this.CheckboxWid.uuids;
        this.table = new SnapshotTableWid({
            id: this.child('#display-table'),
            inst: this.inst,
        });

        // register button's click.
        $(this.child('#remove')).on("click", (e) => {
            let uuids = this.uuids.store.slice();
            new ConfirmActionWid({
                id: this.confirm,
                action: "remove",
                name: uuids.join(", "),
                message: "remove",
            }).onsubmit(() => {
                new SnapshotApi({
                    inst: this.inst,
                    uuids: uuids,
                    name: this.name}).delete();
            });
            $(this.confirm).modal("show");
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
                this.CheckboxWid.refresh();
            });
        });
        this.table.refresh((e) => {
            this.CheckboxWid.refresh();
        });
    }

    create(data) {
        new SnapshotApi({inst: this.inst, name: this.name}).create(data);
    }
}
