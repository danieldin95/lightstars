import {LeaseApi} from "../api/lease.js";
import {LeaseTable} from "../widget/lease/table.js";
import {CheckBoxTab} from "../widget/checkbox/checkbox.js";


class CheckBox extends CheckBoxTab {
}


export class LeasesCtl {
    // {
    //   id: '#network #leases',
    //   uuid: uuid of network,
    //   name: name of network,
    // }
    constructor(props) {
        this.id = props.id;
        this.name = props.name;
        this.uuid = props.uuid;

        this.checkbox = new CheckBox(props);
        this.uuids = this.checkbox.uuids;
        this.table = new LeaseTable({
            id: `${this.id} #display-table`,
            uuid: this.uuid,
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
}
