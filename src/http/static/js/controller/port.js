import {Controller} from "./controller.js";
import {PortTable} from "../widget/port/table.js";
import {CheckBox} from "../widget/common/checkbox.js";


class CheckBoxCtl extends CheckBox {
}


export class PortCtl extends Controller {
    // {
    //   id: '#network #port',
    //   bridge: bridge of network,
    //   name: name of network,
    // }
    constructor(props) {
        super(props);
        this.name = props.name;
        this.uuid = props.uuid;

        this.checkbox = new CheckBoxCtl(props);
        this.uuids = this.checkbox.uuids;
        this.table = new PortTable({
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
