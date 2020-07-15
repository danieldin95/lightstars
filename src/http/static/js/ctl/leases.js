import {Ctl} from "./ctl.js";
import {LeaseApi} from "../api/lease.js";
import {LeaseTable} from "../widget/lease/table.js";
import {CheckBox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckBox {
}


export class LeasesCtl extends Ctl {
    // {
    //   id: '#network #leases',
    //   uuid: uuid of network,
    //   name: name of network,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.uuid = props.uuid;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new LeaseTable({
            id: this.child('#display-table'),
            uuid: this.uuid,
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
}
