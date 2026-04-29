import {Controller} from "./controller.js";
import {LeaseTableWid} from "../widget/lease/leasetable.js";
import {CheckboxWid} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckboxWid {
}


export class LeasesCtl extends Controller {
    // {
    //   id: '#network #leases',
    //   uuid: uuid of network,
    //   name: name of network,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.uuid = props.uuid;

        this.CheckboxWid = new CheckBoxCtl(props);
        this.uuids = this.CheckboxWid.uuids;
        this.table = new LeaseTableWid({
            id: this.child('#display-table'),
            uuid: this.uuid,
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
}
