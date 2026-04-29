import {Controller} from "./controller.js";
import {leasetable} from "../widget/lease/leasetable.js";
import {checkbox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends checkbox {
}


export class Leases extends Controller {
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
        this.table = new leasetable({
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
